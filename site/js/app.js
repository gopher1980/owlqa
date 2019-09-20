function new_graph(json){
var nodes = [];
var edges = [];
var i = 0
var j = 0
var calculateGrpah = function(test, group, last_id) {
        group=group+1
        var dashes = false
        if (last_id != undefined) {
            dashes = true
        }
        
        Object.keys(test).sort().forEach(function(key) {
            var id=id+"_"+i
            nodes[i] = {id:id, label:key, group: group}       
            if (last_id != undefined){
                edges[j] = {from: last_id, to: id, dashes:dashes}
                dashes=false
                j++
            }
            last_id = id
            i++
            if (test[key].fork != undefined){
                calculateGrpah(test[key].fork ,group,id)
            }
        });
    }
    calculateGrpah(json,0)
    return {nodes:nodes, edges:edges}
}

function draw(json_code){
    // create a network
    var len = undefined;

    var test ={}
    try{
        test = JSON.parse(json_code)
    }catch(e){
        test={}
    }
    var container = document.getElementById('mynetwork');
    var data=new_graph(test)
    var options = {
        edges: {
            smooth: {
                type: 'cubicBezier',
                forceDirection: 'horizontal',
                roundness: 0.4
            }
        },
        layout: {
            hierarchical: {
                direction: "UD"
            }
        },
        physics:false
    };
    network = new vis.Network(container, data, options);
}
draw(editor.getValue())