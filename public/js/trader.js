'use strict'

requirejs.config({
  baseUrl: '/public/js',
  paths: {
    moment: 'moment.min',
    d3: 'd3.min',
    dr: 'dom-ready'
  },
})

requirejs(['moment', 'd3'])
requirejs([
  'd3',
  'trader-graph',
  'trader-socket'
], (d3) => {
  console.log(`Successfully loaded D3 version ${d3.version}`)
})
