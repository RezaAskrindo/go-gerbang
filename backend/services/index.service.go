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
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Gerbang - API Gateway</title>
    <style>
        :root {
          --text: #6b6375;
          --text-h: #08060d;
          --bg: #fff;
          --border: #e5e4e7;
          --code-bg: #f4f3ec;
          --accent: #84fc96;
					--accent-bg: rgba(132, 252, 146, 0.15);
					--accent-border: rgba(132, 252, 134, 0.5);
          --social-bg: rgba(244, 243, 236, 0.5);
          --shadow:
            rgba(0, 0, 0, 0.1) 0 10px 15px -3px, rgba(0, 0, 0, 0.05) 0 4px 6px -2px;

          --sans: system-ui, 'Segoe UI', Roboto, sans-serif;
          --heading: system-ui, 'Segoe UI', Roboto, sans-serif;
          --mono: ui-monospace, Consolas, monospace;

          font: 18px/145%% var(--sans);
          letter-spacing: 0.18px;
          color-scheme: light dark;
          color: var(--text);
          background: var(--bg);
          font-synthesis: none;
          text-rendering: optimizeLegibility;
          -webkit-font-smoothing: antialiased;
          -moz-osx-font-smoothing: grayscale;
        }

        @media (max-width: 1024px) {
          :root {
            font-size: 16px;
          }
        }

        @media (prefers-color-scheme: dark) {
          :root {
            --text: #9ca3af;
            --text-h: #f3f4f6;
            --bg: #16171d;
            --border: #2e303a;
            --code-bg: #1f2028;
            --accent: #84fc96;
            --accent-bg: rgba(132, 252, 146, 0.15);
            --accent-border: rgba(132, 252, 134, 0.5);
            --social-bg: rgba(47, 48, 58, 0.5);
            --shadow:
              rgba(0, 0, 0, 0.4) 0 10px 15px -3px, rgba(0, 0, 0, 0.25) 0 4px 6px -2px;
          }
        }

        * {
          margin: 0;
          padding: 0;
          box-sizing: border-box;
        }

        body {
          margin: 0;
          padding: 0;
        }

        h1, h2 {
          font-family: var(--heading);
          font-weight: 500;
          color: var(--text-h);
        }

        h1 {
          font-size: 56px;
          letter-spacing: -1.68px;
          margin: 32px 0 16px;
        }

        h2 {
          font-size: 24px;
          line-height: 118%%;
          letter-spacing: -0.24px;
          margin: 0 0 8px;
        }

        p {
          margin: 8px 0;
          font-size: 18px;
        }

        code {
          font-family: var(--mono);
          font-size: 15px;
          line-height: 135%%;
          padding: 4px 8px;
          background: var(--code-bg);
          border-radius: 4px;
          color: var(--text-h);
          display: inline-block;
        }

        #app {
          width: 1126px;
          max-width: 100%%;
          margin: 0 auto;
          text-align: center;
          border-inline: 1px solid var(--border);
          min-height: 100svh;
          display: flex;
          flex-direction: column;
          box-sizing: border-box;
        }

        #center {
          display: flex;
          flex-direction: column;
          gap: 25px;
          place-content: center;
          place-items: center;
          flex-grow: 1;
          padding: 32px 20px;
        }

        .info-container {
          background: var(--accent-bg);
          border: transparent;
          border-radius: 8px;
          padding: 24px;
          margin-top: 16px;
          display: flex;
          flex-direction: column;
          gap: 12px;
          min-width: 300px;
        }

        .info-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 8px 0;
          border-bottom: 1px solid var(--border);
        }

        .info-item:last-child {
          border-bottom: none;
        }

        .info-label {
          font-weight: 600;
          color: var(--text-h);
          text-align: left;
          flex: 0 0 auto;
        }

        .info-value {
          color: var(--accent);
          font-family: var(--mono);
          text-align: right;
          flex: 1;
          margin-left: 16px;
        }

        @media (max-width: 1024px) {
          h1 {
            font-size: 36px;
            margin: 20px 0 12px;
          }

          #center {
            padding: 24px 16px;
          }

          .info-container {
            min-width: auto;
            width: 100%%;
          }
        }
    </style>
</head>
<body>
    <div id="app">
        <div id="center">
            <h1>Hi There</h1>
            
            <div class="info-container">
                <div class="info-item">
                    <span class="info-label">Name:</span>
                    <span class="info-value">%s</span>
                </div>
                <div class="info-item">
                    <span class="info-label">Version:</span>
                    <span class="info-value">%s</span>
                </div>
                <div class="info-item">
                    <span class="info-label">Maker:</span>
                    <span class="info-value">%s</span>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`, info.Name, info.Version, info.Maker)

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(html)
}
