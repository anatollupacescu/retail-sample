package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("» ")
}

var port = flag.Int("port", 3000, "port to listen to")

func main() {
	http.HandleFunc("/", indexFunc)
	http.HandleFunc("/data.json", dataFunc)

	portStr := fmt.Sprintf(":%d", *port)

	log.Printf("listening on port %s", portStr)

	if err := http.ListenAndServe(portStr, nil); err != nil {
		panic(err)
	}
}

var data = `{
	"nodes": [{"id": "No tests ran yet", "status": "pending"}],
	"links": []
}`

func dataFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, data)
		return
	}

	bts, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data = string(bts)
}

func indexFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, index)
}

var index = `<!DOCTYPE html>
<meta charset="utf-8" />
<style>
  .links line {
    stroke: #999;
    stroke-opacity: 0.6;
  }

  .nodes circle {
    stroke: #fff;
    stroke-width: 1.5px;
  }

  text {
    font-family: sans-serif;
    font-size: 15px;
  }
</style>
<svg width="960" height="600"></svg>
<script src="https://d3js.org/d3.v4.js"></script>
<script>
  var svg = d3.select("svg"),
    width = +svg.attr("width"),
    height = +svg.attr("height");

  var color = ['gray', 'red', 'green'];

  var simulation = d3
    .forceSimulation()
    .force(
      "link",
      d3
        .forceLink()
        .id(function (d) {
          return d.id;
        })
        .distance(190)
    )
    .force("charge", d3.forceManyBody())
    .force("center", d3.forceCenter(width / 2, height / 2));

  d3.json("/data.json", function (error, graph) {
    if (error) throw error;

    svg.append("svg:defs").selectAll("marker")
    .data(["end"])      // Different link/path types can be defined here
    .enter().append("svg:marker")    // This section adds in the arrows
      .attr("id", String)
      .attr("viewBox", "0 -5 10 10")
      .attr("refX", 21)
      .attr("refY", 0.5)
      .attr("markerWidth", 10)
      .attr("markerHeight", 10)
      .attr("orient", "auto")
    .append("svg:path")
      .attr("d", "M0,-5L10,0L0,5");
      
    var link = svg
      .append("g")
      .attr("class", "links")
      .selectAll("line")
      .data(graph.links)
      .enter()
      .append("line")
      .attr("marker-end", "url(#end)")
      .attr("stroke-width", function (d) {
        return Math.sqrt(d.value);
      });

    var node = svg
      .append("g")
      .attr("class", "nodes")
      .selectAll("g")
      .data(graph.nodes)
      .enter()
      .append("g");

    var circles = node
      .append("circle")
      .attr("r", 20)
      .attr("fill", function (d) {
        return color[d.group];
      })
      .call(
        d3
          .drag()
          .on("start", dragstarted)
          .on("drag", dragged)
          .on("end", dragended)
      );

    var lables = node
      .append("text")
      .text(function (d) {
        return d.id + " (" + d.status + ")";
      })
      .attr("x", 6)
      .attr("y", 3);

    node.append("title").text(function (d) {
      return d.id;
    });

    simulation.nodes(graph.nodes).on("tick", ticked);

    simulation.force("link").links(graph.links);

    function ticked() {
      link
        .attr("x1", function (d) {
          return d.source.x;
        })
        .attr("y1", function (d) {
          return d.source.y;
        })
        .attr("x2", function (d) {
          return d.target.x;
        })
        .attr("y2", function (d) {
          return d.target.y;
        });

      node.attr("transform", function (d) {
        return "translate(" + d.x + "," + d.y + ")";
      });
    }
  });

  function dragstarted(d) {
    if (!d3.event.active) simulation.alphaTarget(0.3).restart();
    d.fx = d.x;
    d.fy = d.y;
  }

  function dragged(d) {
    d.fx = d3.event.x;
    d.fy = d3.event.y;
  }

  function dragended(d) {
    if (!d3.event.active) simulation.alphaTarget(0);
    d.fx = null;
    d.fy = null;
  }
</script>
`
