$(document).ready(function() {
  var t = $("#inventoryTable").DataTable({
    ajax: "http://127.0.0.1:8080/inbound/config",
    columns: [{ data: "id" }, { data: "type" }]
  });

  $("#inventoryTable tbody").on("click", "tr", function() {
    $(this).toggleClass("list-group-item-dark");
  });

  var nameInput = $("#name");

  nameInput.keyup(function() {
    $(".invalid-feedback").removeClass("d-block");
  });

  var form = $("#mainForm");

  form.on("submit", function(e) {
    e.preventDefault();

    var data = nameInput.val();

    if (!data) {
      $(".invalid-feedback").addClass("d-block");
      return;
    }

    $.ajax({
      type: "POST",
      url: "/inbound/config",
      data: JSON.stringify([data]),
      contentType: "application/json",
      success: function() {
        t.ajax.reload();
        nameInput.val("");
      },
      error: function(resp) {
        console.log(resp.statusText);
      }
    });
  });
});
