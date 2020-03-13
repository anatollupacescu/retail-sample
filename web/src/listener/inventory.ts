import $ = require("jquery");

import RetailInventory from "../retailapp/inventory";

export function initializeInventory(app: RetailInventory) {
  app.fetchInventoryState().then(() => {
    populateTable(app.getInventory());
  });

  $("#inventoryTable tbody").on("click", "tr", function() {
    $(this).toggleClass("list-group-item-dark");
  });

  let nameInput = $("#name");

  nameInput.keyup(function() {
    $("#nonempty.invalid-feedback").removeClass("d-block");
    $("#unique.invalid-feedback").removeClass("d-block");
  });

  $("#mainForm").on("submit", function(e) {
    e.preventDefault();

    var data = <string>nameInput.val();

    app
      .addInventoryItem(data)
      .then(rsp => {
        if (rsp === "name empty") {
          $("#nonempty.invalid-feedback").addClass("d-block");
          return;
        }
        if (rsp === "name present") {
          $("#unique.invalid-feedback").addClass("d-block");
          return;
        }
        $("#inventoryTable tbody tr").remove();
        populateTable(app.getInventory());
        nameInput.val("");
      })
      .catch(err => {
        if (err === "ERR_EMPTY") {
          $("#nonempty.invalid-feedback").addClass("d-block");
        }
        if (err === "ERR_UNIQUE") {
          $("#unique.invalid-feedback").addClass("d-block");
        }
      });
  });
}

interface inventoryItem {
  id: number;
  name: string;
}

function populateTable(items: inventoryItem[]) {
  let table = $("#inventoryTable tbody")[0];
  let rows = items.sort((i1: inventoryItem, i2: inventoryItem) => {
    return i1.id - i2.id;
  });
  rows.forEach((element: inventoryItem) => {
    let row = table.insertRow(0);
    row.insertCell(0).innerHTML = element.id;
    row.insertCell(1).innerHTML = element.name;
  });
}
