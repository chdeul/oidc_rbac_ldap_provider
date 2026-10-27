    // Fragments count in progress
    var nbEmbed = 0

    class HTTPError extends Error {
        constructor(code, message, extras) {
          super(message);
          this.name = message;
          this.statusCode = code;
        }
    }

    // Load pages fragments <v-embed> .vue
    customElements.define('v-embed', class extends HTMLElement {
        async connectedCallback() {
            nbEmbed++
            let src = this.getAttribute('src')
            this.outerHTML = await (await fetch(src)).text()       
            nbEmbed--
        }
    })

    // Wait ends of loading...
    var tryInit = 1 // 8 Seconds
    function loadingVue() {

        // Wait all fetch (v-embed) done, 4s max
        while ( nbEmbed > 0 ) { 
            tryInit--
            if ( tryInit > 0 )
                setTimeout( loadingVue, 200 ) // Milliseconds   
            else
                throw new HTTPError(404,"Http timeout at page-fragments loading"); 
            return
        }

        // Eval scripts from v-embed
        var scripts = document.body.querySelectorAll("script")
        for (var i = 0; i < scripts.length; i++) 
            window.eval(scripts[i].innerText) 
    }

    window.onload = function () {

        // Wait loading...
        loadingVue()
    }
