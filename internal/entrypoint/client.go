package entrypoint

var clientRenderFunction = `hydrateRoot(document.getElementById("root"), <App {...props} />);`
var clientRenderFunctionWithLayout = `hydrateRoot(document.getElementById("root"), <Layout><App {...props} /></Layout>);`
