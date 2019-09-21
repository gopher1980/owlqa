function pagination(current) {
    page = current
    var p1 = parseInt(current)
    var p0 = p1 - 1
    if (p0 <= 0) p0 = 1;
    var p2 = p1 + 1
    var p3 = p1 + 2
    $("#pagination").html(
        `<li class="page-item disabled"><a href="?page=${p0}">Previous</a></li>
          <li class="page-item active"><a href="?page=${p1}" class="page-link">${p1}</a></li>
          <li class="page-item"><a href="?page=${p2}" class="page-link">${p2}</a></li>
          <li class="page-item "><a href="?page=${p3}" class="page-link">${p3}</a></li>
          <li class="page-item"><a href="?page=${p2}" class="page-link">Next</a></li>
    `
    )
}


function delete_show(id) {
    globalId = id
    $("#id").html(id)
    $("#delete").modal("show")
}

function real_delete(refresh) {
    $("#delete").modal("hide")
    $.ajax({
        url: url + "/" + globalId,
        type: "DELETE",
        success: function(data){
            refresh()
        }
    })
    globalId = -1
    

}