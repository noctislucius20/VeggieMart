// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs'

export default withNuxt(
  // Your custom configs here
  {
    rules: {
      // Mengizinkan elemen void HTML (seperti <hr/> atau <br/>) menggunakan self-closing atau tidak ("any")
      'vue/html-self-closing': [
        'error',
        {
          html: {
            void: 'any', // Ubah dari "never" menjadi "any" atau "always"
            normal: 'always',
            component: 'always'
          }
        }
      ]
    }
  }
)
