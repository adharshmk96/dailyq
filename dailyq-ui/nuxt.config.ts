// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxt/ui'],
  css: ['~/assets/css/main.css'],

  // The build is a static SPA so it can be embedded into the Go binary and
  // served by the API server itself.
  ssr: false,

  runtimeConfig: {
    public: {
      // Same-origin by default: the Go server serves both the SPA and the API.
      // In dev, nitro.devProxy below forwards /api to the local API server.
      apiBase: '/api/v1'
    }
  },

  nitro: {
    devProxy: {
      '/api': {
        target: process.env.NUXT_DEV_API_TARGET || 'http://localhost:8080/api',
        changeOrigin: true
      }
    }
  },

  app: {
    head: {
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon/favicon-32x32.png' },
        { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon/favicon-16x16.png' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/favicon/apple-touch-icon.png' },
        { rel: 'manifest', href: '/favicon/site.webmanifest' }
      ],
      meta: [
        { name: 'theme-color', content: '#ffffff' }
      ]
    }
  },

  icon: {
    provider: 'none',
    clientBundle: {
      scan: true
    }
  }
})
