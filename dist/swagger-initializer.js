window.onload = function () {
    //<editor-fold desc="Changeable Configuration Block">

    function excludeServiceName(name) {
        const prefixes = ["grpc.", "envoy.", "xds."];

        for (let prefix of prefixes) {
            if (name.startsWith(prefix)) {
                return false;
            }
        }
        return true;
    }
    var currentPath = window.location.pathname;
    // 使用正则表达式匹配路径，剔除包含 'swagger' 的部分
    var apiPrefix = currentPath.replace(/\/swagger\/?.*$/, '');
    let servicesUrl = new URL("/swagger-query/services", window.location.href);
    if (apiPrefix.length > 0) {
        console.log("当前页面的 API 前缀为: " + apiPrefix);
        servicesUrl = new URL(apiPrefix + "/swagger-query/services", window.location.href);
    }

    fetch(servicesUrl.toString())
        .then(response => response.json())
        .then(data => {
            const urls = data.services.filter(x => excludeServiceName(x)).map((x) => {
                let url = new URL("/swagger-query/service/" + x, window.location.href);
                if (apiPrefix.length > 0) {
                    url = new URL(apiPrefix + "/swagger-query/service/" + x, window.location.href);
                }
                return {url: url.toString(), name: x}
            });
            console.log(urls)
            // Begin Swagger UI call region
            const ui = SwaggerUIBundle({
                urls: urls,
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset,
                ],
                plugins: [
                    SwaggerUIBundle.plugins.Topbar,
                    SwaggerUIBundle.plugins.DownloadUrl,
                ],
                layout: "StandaloneLayout"
            });
            // End Swagger UI call region

            window.ui = ui;
        });
    // the following lines will be replaced by docker/configurator, when it runs in a docker-container
    // window.ui = SwaggerUIBundle({
    //   url: "https://petstore.swagger.io/v2/swagger.json",
    //   // urls: ["a", "b"],
    //   dom_id: '#swagger-ui',
    //   deepLinking: true,
    //   presets: [
    //     SwaggerUIBundle.presets.apis,
    //     SwaggerUIStandalonePreset
    //   ],
    //   plugins: [
    //     SwaggerUIBundle.plugins.DownloadUrl
    //   ],
    //   layout: "StandaloneLayout"
    // });

    //</editor-fold>
};
