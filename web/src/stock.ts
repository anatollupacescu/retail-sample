$(document).ready(function() {
  var stockTable = $("#stock").DataTable({
    ajax: "http://127.0.0.1:8080/stock",
    columns: [{ data: "name" }, { data: "qty", searchable: false }]
  });

  var dropDown = $("#inventoryItemType");

  $.ajax({
    type: "GET",
    url: "/inbound/config",
    accept: "application/json",
    success: function(data) {
      let items = data.data;
      if (items) {
        dropDown.empty();
        items.map(function(item) {
          dropDown.append(new Option(item.type));
        });
      }
    },
    error: function(resp) {
      console.log(resp.statusText);
    }
  });

  $("#add").on("click", function() {
    var name = dropDown.val();
    var qty = $("#count").val();
    if (!name || !qty) {
      console.log("name and quantity mandatory");
      return;
    }
    var payload = {};
    payload[name] = Number.parseInt(qty);
    $.ajax({
      type: "POST",
      url: "/inbound",
      data: JSON.stringify(payload),
      contentType: "application/json",
      success: function() {
        stockTable.ajax.reload();
      },
      error: function(resp) {
        console.log(resp.statusText);
      }
    });
  });
});
