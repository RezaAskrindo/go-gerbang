# CATATAN

## Penggunaan S3 Pada Caddy Server

1. base di ganti jadi folder misal "/mfe/react-app-shell"
2. contoh config di JSON API Caddy:
**{
  "match": [
    {
      "path": [
        "/mfe/react-app-shell/assets/*",
        "/mfe/react-app-shell/vite.svg"
      ]
    }
  ],
  "handle": [
    {
      "handler": "reverse_proxy",
      "upstreams": [
        {
          "dial": "s3.nevaobjects.id:443"
        }
      ],
      "transport": {
        "protocol": "http",
        "tls": {}
      },
      "headers": {
        "request": {
          "set": {
            "Host": [
              "s3.nevaobjects.id"
            ]
          }
        }
      }
    }
  ]
},
{
  "match": [
    {
      "host": ["auth.siskor.web.id"]
    }
  ],
  "handle": [
    {
      "handler": "rewrite",
      "uri": "/mfe/react-app-shell/index.html"
    },
    {
      "handler": "reverse_proxy",
      "upstreams": [
        {
          "dial": "s3.nevaobjects.id:443"
        }
      ],
      "transport": {
        "protocol": "http",
        "tls": {}
      },
      "headers": {
        "request": {
          "set": {
            "Host": [
              "s3.nevaobjects.id"
            ]
          }
        }
      }
    }
  ]
}**