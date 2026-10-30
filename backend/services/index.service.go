package services

import (
	"fmt"

	"go-gerbang/config"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v3"
)

func MainService(c fiber.Ctx) error {
	info := types.ApiInfo{
		Name:    config.AppName,
		Version: config.VesionApp,
		Maker:   "Muhammad Reza",
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="color-scheme" content="dark">
    <title>%s - API Gateway</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        :root {
            --bg: #000000;
            --bg-rgb: 0 0 0;
            --text: #ffffff;
            --text-muted: #a1a1aa;
            --outline: #128052;
            --accent: #128052;
            --accent-rgb: 50 255 170;
            --border: #1f1f1f;
            --code-bg: #0a0a0a;

            --hero-fade-mid: 0.4;
            --hero-fade-end: 100%%;

            --sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Inter', 'Helvetica Neue', sans-serif;
            --mono: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
            color-scheme: dark;
        }

        html, body {
            min-height: 100vh;
            background: var(--bg);
            color: var(--text);
            font-family: var(--sans);
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
        }

        body {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 64px 24px;
            text-align: center;
            overflow-x: hidden;
        }

        .hero {
            width: 100%%;
            max-width: 1100px;
            display: flex;
            flex-direction: column;
            align-items: center;
        }

        .brand-wrap {
            position: relative;
            width: 100%%;
            max-width: 1000px;
            pointer-events: none;
            user-select: none;
        }

        .brand-outline {
            display: block;
            width: 100%%;
            height: auto;
            overflow: visible;
            color: var(--outline);
        }

        .brand-fade {
            position: absolute;
            top: 0;
            bottom: -10px;
            left: -16px;
            right: -16px;
            pointer-events: none;
            background: linear-gradient(
                to bottom in srgb,
                rgb(var(--bg-rgb) / 0) 0%%,
                rgb(var(--bg-rgb) / var(--hero-fade-mid)) 50%%,
                rgb(var(--bg-rgb) / 1) var(--hero-fade-end)
            );
        }

        /* ---- shimmer line ---- */
        .brand-shimmer {
            position: absolute;
            left: 50%%;
            transform: translateX(-50%%);
            width: min(70%%, 720px);
            height: 1px;
            border-radius: 999px;
            overflow: hidden;
            pointer-events: none;
            z-index: 2;
            background: linear-gradient(
                90deg,
                transparent 0%%,
                rgb(var(--accent-rgb) / 0.35) 30%%,
                rgb(var(--accent-rgb) / 0.65) 50%%,
                rgb(var(--accent-rgb) / 0.35) 70%%,
                transparent 100%%
            );
            box-shadow:
                0 0 10px rgb(var(--accent-rgb) / 0.35),
                0 0 24px rgb(var(--accent-rgb) / 0.15);
        }

        .brand-shimmer::after {
            content: '';
            position: absolute;
            inset: 0;
            background: linear-gradient(
                90deg,
                transparent 0%%,
                transparent 35%%,
                rgba(255, 255, 255, 0.95) 50%%,
                transparent 65%%,
                transparent 100%%
            );
            background-size: 200%% 100%%;
            animation: brand-shimmer 3s ease-in-out infinite;
        }

        @keyframes brand-shimmer {
            0%%   { background-position: 200%% 0; }
            100%% { background-position: -200%% 0; }
        }

        @media (prefers-reduced-motion: reduce) {
            .brand-shimmer::after { animation: none; }
        }
        /* ---------------------- */

        h1 {
            margin-top: 5px;
            font-size: clamp(26px, 5vw, 56px);
            font-weight: 600;
            line-height: 1.08;
            letter-spacing: -0.03em;
            margin-bottom: 20px;
            max-width: 760px;
        }

        .subtitle {
            font-size: clamp(15px, 1.6vw, 18px);
            color: var(--text-muted);
            max-width: 580px;
            line-height: 1.55;
            margin-bottom: 40px;
        }

        .subtitle strong {
            color: var(--text);
            font-weight: 600;
        }

        .command-box {
            display: inline-flex;
            align-items: center;
            gap: 10px;
            background: var(--code-bg);
            border: 1px solid var(--border);
            border-radius: 999px;
            padding: 12px 20px;
            font-family: var(--mono);
            font-size: 14px;
            margin-bottom: 56px;
            transition: border-color 0.2s ease;
        }

        .command-box:hover {
            border-color: #333;
        }

        .command-box .prompt {
            color: #555;
            user-select: none;
        }

        .command-box .cmd {
            color: var(--text);
        }

        .meta {
            display: flex;
            flex-wrap: wrap;
            justify-content: center;
            align-items: center;
            gap: 20px;
            font-size: 14px;
            color: var(--text-muted);
        }

        .meta-item {
            display: inline-flex;
            align-items: center;
            gap: 8px;
        }

        .meta-label {
            color: #666;
        }

        .meta-value {
            color: var(--text);
            font-family: var(--mono);
            font-weight: 500;
        }

        .meta-value.accent {
            color: var(--accent);
        }

        .divider {
            width: 1px;
            height: 14px;
            background: var(--border);
        }

        @media (max-width: 640px) {
            body { padding: 48px 20px; }
            .brand-fade { bottom: -24px; }
            .brand-shimmer {
                bottom: -5px;
                width: 78%%;
            }
            h1 { margin-top: 52px; }
            .command-box {
                font-size: 12px;
                padding: 10px 16px;
                margin-bottom: 40px;
            }
            .meta { gap: 12px; font-size: 13px; }
            .divider { display: none; }
        }
    </style>
</head>
<body>
    <main class="hero">

        <div class="brand-wrap">
            <svg
                class="brand-outline"
                viewBox="0 0 1300 200"
                preserveAspectRatio="xMidYMid meet"
                xmlns="http://www.w3.org/2000/svg"
                aria-hidden="true"
                focusable="false">
                <text
                    x="650"
                    y="172"
                    text-anchor="middle"
                    font-family="ui-sans-serif, system-ui, -apple-system, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif"
                    font-weight="900"
                    font-size="230"
                    letter-spacing="2"
                    textLength="1280"
                    lengthAdjust="spacingAndGlyphs"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="3"
                    stroke-linejoin="round"
                    stroke-linecap="round">GERBANG</text>
            </svg>
            <div class="brand-fade" aria-hidden="true"></div>
            <div class="brand-shimmer" aria-hidden="true"></div>
        </div>

        <h1>Fastest API Gateway using Gofiber</h1>

        <p class="subtitle">
            Used by developers who need speed. <strong>Gerbang</strong> is a lightweight
            API gateway built on <strong>Gofiber</strong> — the fastest HTTP engine for Go.
        </p>

        <div class="command-box">
            <span class="prompt">$</span>
            <span class="cmd">gerbang start</span>
        </div>

        <div class="meta">
            <div class="meta-item">
                <span class="meta-label">Version</span>
                <span class="meta-value accent">%s</span>
            </div>
            <div class="divider"></div>
            <div class="meta-item">
                <span class="meta-label">Created by</span>
                <span class="meta-value">%s</span>
            </div>
        </div>
    </main>
</body>
</html>`, info.Name, info.Version, info.Maker)

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(html)
}
