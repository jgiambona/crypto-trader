'use strict'

requirejs.config({
  baseUrl: '/public/js',
  paths: {
    moment: 'moment.min',
    dr: 'dom-ready'
  },
})

requirejs(['moment'])
requirejs([
], () => {
  window.setInterval(() => {
    console.log("reload frames")
    document.frames["price-graph"].location.reload()
    document.frames["win-loss-graph"].location.reload()
  }, 3000)
})
