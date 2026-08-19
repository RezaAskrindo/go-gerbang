export interface FlatCaddyConfig {
  i?: number;
  type?: string;
  rewrite_uri?: string;
  rewrite_strip_path_prefix?: string;
  file_server?: string;
  reverse_proxy?: string;
  match_host?: string[];
  match_path?: string[];
  host?: string;
}

interface CaddyMatch {
  host?: string[];
  path?: string[];
}

type CaddyHandler =
  | { handler: "rewrite"; uri?: string; strip_path_prefix?: string }
  | {
      handler: "reverse_proxy";
      upstreams: { dial: string }[];
      transport?: { protocol: string; tls: Record<string, never> };
      headers?: { request: { set: { Host: string[] } } };
    }
  | { handler: "file_server"; root: string };

interface CaddyRoute {
  match?: CaddyMatch[];
  handle: CaddyHandler[];
}

export interface CaddyConfig {
  listen: string[];
  routes: CaddyRoute[];
}

// A "TLS reverse_proxy to a fixed host" handler is the signature of the
// S3-backed FE pattern (as opposed to a plain dial-based reverse_proxy).
function extractS3Host(route: CaddyRoute): string | undefined {
  const handler = route.handle?.find(
    (h): h is Extract<CaddyHandler, { handler: "reverse_proxy" }> =>
      h.handler === "reverse_proxy" && !!h.transport?.tls
  );
  if (!handler) return undefined;
  return handler.headers?.request?.set?.Host?.[0] ?? handler.upstreams?.[0]?.dial;
}

function hasRewrite(route: CaddyRoute): boolean {
  return !!route.handle?.some((h) => h.handler === "rewrite");
}

function hasMatch(route: CaddyRoute): boolean {
  return !!(route.match && route.match.length);
}

function mapSingleRoute(route: CaddyRoute, index: number): FlatCaddyConfig {
  const result: FlatCaddyConfig = {
    i: index,
    rewrite_uri: undefined,
    rewrite_strip_path_prefix: undefined,
    file_server: undefined,
    reverse_proxy: undefined,
    match_host: undefined,
    match_path: undefined,
    host: undefined,
  };

  if (Array.isArray(route.handle)) {
    for (const handler of route.handle) {
      if (handler.handler === "rewrite") {
        if (handler.uri) result.rewrite_uri = handler.uri;
        if (handler.strip_path_prefix) result.rewrite_strip_path_prefix = handler.strip_path_prefix;
      } else if (handler.handler === "file_server") {
        result.file_server = handler.root;
      } else if (handler.handler === "reverse_proxy") {
        if (handler.transport?.tls) {
          // S3-style TLS proxy - capture as `host`, not `reverse_proxy`
          result.host = extractS3Host(route);
        } else {
          result.reverse_proxy = handler.upstreams?.[0]?.dial;
        }
      }
    }
  }

  if (Array.isArray(route.match)) {
    for (const matchRule of route.match) {
      // BUG FIX: keep the full array - the old code did matchRule.host[0] /
      // matchRule.path[0], silently dropping every entry after the first.
      if (matchRule.host) result.match_host = matchRule.host;
      if (matchRule.path) result.match_path = matchRule.path;
    }
  }

  return result;
}

export function transformCaddyConfig(config: CaddyConfig): FlatCaddyConfig[] {
  const result: FlatCaddyConfig[] = [];
  const routes = config?.routes;

  if (!routes) {
    return result;
  }

  for (let idx = 0; idx < routes.length; idx++) {
    const route = routes[idx];
    const host = extractS3Host(route);

    // Look for the exact pair reverseTransformCaddyConfig produces for an S3
    // app: [match+reverse_proxy route] immediately followed by
    // [rewrite+reverse_proxy route, no match], both hitting the same host.
    // Merge them back into ONE flat row instead of two.
    if (host && hasMatch(route) && !hasRewrite(route)) {
      const next = routes[idx + 1];
      const nextHost = next ? extractS3Host(next) : undefined;

      if (
        next &&
        nextHost === host &&
        hasRewrite(next) &&
        (!hasMatch(next) || (next.match?.length === 1 && next.match[0].host && !next.match[0].path))
      ) {
        const matchRule = route.match![0];
        const rewriteHandler = next.handle.find(
          (h): h is Extract<CaddyHandler, { handler: "rewrite" }> => h.handler === "rewrite"
        );

        result.push({
          i: result.length,
          host,
          match_host: matchRule.host,
          match_path: matchRule.path,
          rewrite_uri: rewriteHandler?.uri,
          rewrite_strip_path_prefix: rewriteHandler?.strip_path_prefix,
          file_server: undefined,
          reverse_proxy: undefined,
        });

        idx++; // consume the paired route too, don't emit it separately
        continue;
      }
    }

    result.push(mapSingleRoute(route, result.length));
  }

  return result;
  // if (!config?.routes) {
  //   return [];
  // }

  // return config.routes.map((route: any, index: number) => {
  //   const result: FlatCaddyConfig = {
  //     i: index,
  //     rewrite_uri: undefined,
  //     rewrite_strip_path_prefix: undefined,
  //     file_server: undefined,
  //     reverse_proxy: undefined,
  //     match_host: undefined,
  //     match_path: undefined
  //   };

  //   // Process handlers
  //   if (route.handle && Array.isArray(route.handle)) {
  //     route.handle.forEach((handler: any) => {
  //       if (handler.handler === 'rewrite') {
  //         if (handler.uri) result.rewrite_uri = handler.uri;
  //         if (handler.strip_path_prefix) result.rewrite_strip_path_prefix = handler.strip_path_prefix;
  //       } else if (handler.handler === 'file_server') {
  //         result.file_server = handler.root;
  //       } else if (handler.handler === 'reverse_proxy') {
  //         result.reverse_proxy = handler.upstreams?.[0]?.dial;
  //       }
  //     });
  //   }

  //   // Process match rules
  //   if (route.match && Array.isArray(route.match)) {
  //     route.match.forEach((matchRule: any) => {
  //       if (matchRule.host) {
  //         result.match_host = matchRule.host[0];
  //       }
  //       if (matchRule.path) {
  //         result.match_path = matchRule.path[0];
  //       }
  //     });
  //   }

  //   return result;
  // });
}

export function reverseTransformCaddyConfig(flatArray: FlatCaddyConfig[], listen: string[] = [":443"]): CaddyConfig {
  const arrayUsed: FlatCaddyConfig[] = flatArray.flatMap((item) => {
    if (!item.host) return [item];

    const rows: FlatCaddyConfig[] = [];

    if (item.match_path?.length || item.match_host?.length) {
      rows.push({ ...item, rewrite_uri: undefined, rewrite_strip_path_prefix: undefined });
    }
    if (item.rewrite_uri) {
      rows.push({ ...item, match_path: undefined, match_host: undefined });
    }

    // host set but neither match_* nor rewrite_uri given - don't silently drop it
    return rows.length ? rows : [item];
  });

  return {
    listen,
    routes: arrayUsed.map((item: FlatCaddyConfig) => {
      const handle: CaddyHandler[] = [];

      // Rewrite handler - truthy checks, not `!== undefined`, so empty-string
      // form defaults ("") don't produce a handler with junk fields.
      if (item.rewrite_uri || item.rewrite_strip_path_prefix) {
        handle.push({
          handler: "rewrite",
          ...(item.rewrite_uri ? { uri: item.rewrite_uri } : {}),
          ...(item.rewrite_strip_path_prefix
            ? { strip_path_prefix: item.rewrite_strip_path_prefix }
            : {}),
        });
      }

      // Plain reverse_proxy, e.g. a backend service dial address
      if (item.reverse_proxy) {
        handle.push({
          handler: "reverse_proxy",
          upstreams: [{ dial: item.reverse_proxy }],
        });
      }

      // Static file server
      if (item.file_server) {
        handle.push({ handler: "file_server", root: item.file_server });
      }

      // S3-backed FE: proxy to the bucket over TLS, forcing the Host header
      if (item.host) {
        handle.push({
          handler: "reverse_proxy",
          upstreams: [{ dial: `${item.host}:443` }],
          transport: { protocol: "http", tls: {} },
          headers: { request: { set: { Host: [item.host] } } },
        });
      }

      const route: CaddyRoute = { handle };

      // host + path go into ONE match object (AND'd), not two separate
      // objects (which Caddy treats as OR'd alternatives).
      if (item.match_host?.length || item.match_path?.length) {
        route.match = [
          {
            ...(item.match_host?.length ? { host: item.match_host } : {}),
            ...(item.match_path?.length ? { path: item.match_path } : {}),
          },
        ];
      }

      return route;
    })
  };
}