$(document).ready(function() {
  let apiUrl = process.env.API_URL;
  var stockTable = $("#stock").DataTable({
    ajax: `${apiUrl}/stock`,
    columns: [{ data: "name" }, { data: "qty", searchable: false }]
  });

  var dropDown = $("#inventoryItemType");

  $.ajax({
    type: "GET",
    url: `${apiUrl}/inbound/config`,
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
      url: `${apiUrl}/inbound`,
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
