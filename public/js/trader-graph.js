'use strict'

define([
  'd3',
  'dr!'
], (d3) => {
    let p = document.getElementById('graph')
    let margin = { top: 20, right: 20, bottom: 30, left: 30 },
      cw = p.clientWidth - 30,
      ch = 500,
      width = cw - margin.left - margin.right,
      height = ch - margin.top - margin.bottom,
      svg = d3.select(p).append('svg:svg')
        .attr('class', 'chart')
        .attr('width', cw)
        .attr('height', ch),
      g = svg.append('g')
        .attr('transform', `translate(${margin.left}, ${margin.top})`),
      legend = svg.append('g')
        .attr('class', 'legend')

    d3.json('/public/data/yahoo_usd_btc.json').then((data) => {
      let meta = data.chart.result[0].meta
      let timestamp = data.chart.result[0].timestamp
      let quote = data.chart.result[0].indicators.quote
      let x = d3.scaleTime()
        .rangeRound([0, width])
      let y = d3.scaleTime()
        .range([height, 0])
      let line = d3.line()
        .x(d => x(d.close))
        .y(d => y(d.close))

      x.domain(d3.extent(timestamp, d => d))
      y.domain(d3.extent(quote, d => d))

      g.append('g')
        .attr('class', 'x-axis')
        .attr('transform', `translate(0, ${height})`)
        .call(d3.axisBottom(x))
        .append('text')
          .attr('class', 'label')
          .attr('x', width)
          .attr('y', -6)
          .style('text-anchor', 'end')
          .text('Month')

      g.append('g')
        .attr('class', 'y-axis')
        .call(d3.axisLeft(y))
        .append('text')
          .attr('fill', '#000')
          .attr('transform', 'rotate(-90)')
          .attr('y', 6)
          .attr('dy', '0.71em')
          .attr('text-anchor', 'end')
          .text('Market Price ($)')

      g.append('path')
        .datum(quote)
        .attr('fill', 'none')
        .attr('stroke', 'steelblue')
        .attr('stroke-linejoin', 'round')
        .attr('stroke-linecap', 'round')
        .attr('stroke-width', 1.5)
        .attr('d', line)
    })
})
