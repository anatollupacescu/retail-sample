$(document).ready(function() {
  var outboundTable = $("#outboundList").DataTable({
    ajax: "http://127.0.0.1:8080/outbound",
    columns: [
      { data: "date" },
      { data: "name" },
      { data: "qty", searchable: false }
    ]
  });

  var dropDown = $("#outboundItemType");
  var reloadItemNameList = function() {
    $.ajax({
      type: "GET",
      url: "/outbound/config",
      accept: "application/json",
      success: function(data) {
        let items = data.data;
        if (items) {
          dropDown.empty();
          items.map(function(item) {
            dropDown.append(new Option(item.name));
          });
        }
      },
      error: function(resp) {
        console.log(resp.statusText);
      }
    });
  };

  reloadItemNameList();

  $("#sell").on("click", function() {
    var name = dropDown.val();
    var qty = $("#count").val();
    if (!name || !qty) {
      console.log("name and quantity mandatory");
      return;
    }
    var payload = {
      name: name,
      qty: Number.parseInt(qty)
    };
    $.ajax({
      type: "POST",
      url: "/outbound",
      data: JSON.stringify(payload),
      contentType: "application/json",
      success: function() {
        outboundTable.ajax.reload();
        reloadItemNameList();
      },
      error: function(resp) {
        console.log(resp.statusText);
      }
    });
  });
});
