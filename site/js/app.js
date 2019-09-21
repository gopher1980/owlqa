function new_graph(json) {
    var nodes = [];
    var edges = [];
    var i = 0
    var j = 0
    var calculateGrpah = function (test, group, last_id) {
        group = group + 1
        var dashes = false
        if (last_id != undefined) {
            dashes = true
        }

        Object.keys(test).sort().forEach(function (key) {
            var id = id + "_" + i
            nodes[i] = {
                id: id,
                label: key,
                group: group
            }
            if (last_id != undefined) {
                edges[j] = {
                    from: last_id,
                    to: id,
                    dashes: dashes
                }
                dashes = false
                j++
            }
            last_id = id
            i++
            if (test[key].fork != undefined) {
                calculateGrpah(test[key].fork, group, id)
            }
        });
    }
    calculateGrpah(json, 0)
    return {
        nodes: nodes,
        edges: edges
    }
}

function draw(json_code) {
    // create a network
    var len = undefined;

    var test = {}
    try {
        test = JSON.parse(json_code)
    } catch (e) {
        test = {}
    }
    var container = document.getElementById('mynetwork');
    var data = new_graph(test)
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
        physics: false
    };
    network = new vis.Network(container, data, options);
}

function formatter() {
    var jsObj = JSON.parse(editor.getValue())
    value = JSON.stringify(jsObj, null, "\t");
    editor.setValue(value)
}



function save(){
    var id = null
    if (globalId != null){
        id = parseInt(globalId)
    }
    
    var element = {
        id: id,
        name : $("#test_name").val(),
        code : editor.getValue()
    }
    $.post("/catalog", JSON.stringify(element)).done(function(data){
        if (data.Error==null){
            window.location.replace("/site/console.html?id="+data.Value.id);
        }else{
            alert(JSON.stringify(data))
        }
    })


   

}

function run() {
    var json = editor.getValue()
    var rand = Math.floor((Math.random() * 10000000) + 1);
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var jsObj = JSON.parse(this.responseText)
            value = JSON.stringify(jsObj, null, "\t");

            editorResult.setValue(value)
            $("#loading").modal("hide")
        }
    };

    $("#loading").modal("show")
    var query = (document.getElementById("query").value).replace("&", "")
    if (query == "") query = "$"
    xhttp.open("POST", "http://127.0.0.1:9090/run?s=SESSION" + rand + "&q=" + query, true);
    xhttp.setRequestHeader("Content-type", "application/json");
    xhttp.send(json);



}



