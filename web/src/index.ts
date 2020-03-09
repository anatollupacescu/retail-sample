import $ from 'jquery'
import axios from 'axios'

function fetchInventoryData (apiUrl: string) {
  axios
    .get(`${apiUrl}/inventory`)
    .then(function (response) {
      populateTable(response.data)
    })
    .catch(function (error) {
      console.error(error)
    })
    .then(function () {
      console.log('done fetching')
    })
}

interface inventoryItem {
  id: number
  name: string
}

function populateTable (req: any) {
  let table = $('#inventoryTable tbody')[0]
  let rows = req.data.sort((i1: inventoryItem, i2: inventoryItem) => {
    return i1.id > i2.id
  })
  rows.forEach((element: inventoryItem) => {
    let row = table.insertRow(0)
    row.insertCell(0).innerHTML = element.id
    row.insertCell(1).innerHTML = element.name
  })
}

$(document).ready(function () {
  let apiUrl = process.env.API_URL

  if (!apiUrl) {
    return
  }

  fetchInventoryData(apiUrl)

  $('#inventoryTable tbody').on('click', 'tr', function () {
    $(this).toggleClass('list-group-item-dark')
  })

  let nameInput = $('#name')

  nameInput.keyup(function () {
    $('#nonempty.invalid-feedback').removeClass('d-block')
    $('#unique.invalid-feedback').removeClass('d-block')
  })

  var form = $('#mainForm')

  form.on('submit', function (e: { preventDefault: () => void }) {
    e.preventDefault()

    var data = nameInput.val()

    if (!data) {
      $('#nonempty.invalid-feedback').addClass('d-block')
      return
    }

    axios
      .post(`${apiUrl}/inventory`, [data])
      .then(function () {
        $('#inventoryTable tbody tr').remove()
        fetchInventoryData(apiUrl)
        nameInput.val('')
      })
      .catch(function (resp) {
        if (resp.response.data === 'ERR_UNIQUE') {
          $('#unique.invalid-feedback').addClass('d-block')
        }
      })
  })
})
