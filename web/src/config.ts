$(document).ready(function() {
  let apiUrl = process.env.API_URL

  var fpTable = $('#finishedProducts').DataTable({
    ajax: `${apiUrl}/outbound/config`,
    columns: [{ data: 'name' }, { data: 'count', searchable: false }]
  })

  var items = $('#items').DataTable({
    columns: [{ data: 'name' }, { data: 'count' }]
  })

  var reloadItemNameList = function() {
    $.ajax({
      type: 'GET',
      url: `${apiUrl}/inbound/config`,
      accept: 'application/json',
      success: function(data) {
        let items = data.data
        if (items) {
          $('#itemType').empty()
          items.map(function(item) {
            $('#itemType').append(new Option(item.type))
          })
        }
      },
      error: function(resp) {
        console.log(resp.statusText)
      }
    })
  }

  reloadItemNameList()

  $('#createNew').on('click', function() {
    var tableData = items.rows().data()
    var fpName = $('#finishedProductName').val()
    console.log('data gt', JSON.stringify(tableData))
    if (!items.rows().count()) {
      $('#noRowsErr.invalid-feedback').addClass('d-block')
    }
    if (!fpName) {
      $('#noNameErr.invalid-feedback').addClass('d-block')
    }
    if (tableData && fpName) {
      var cs = {}
      tableData.each(function(i) {
        cs[i.name] = Number.parseInt(i.count)
      })
      var payload = {
        name: fpName,
        items: cs
      }
      $.ajax({
        type: 'POST',
        url: `${apiUrl}/outbound/config`,
        data: JSON.stringify(payload),
        contentType: 'application/json',
        success: function() {
          items.clear().draw()
          fpTable.ajax.reload()
          $('#finishedProductName').val('')
          reloadItemNameList()
        },
        error: function(resp) {
          if (resp.responseText === 'ERR_UNIQUE') {
            $('#unique.invalid-feedback').addClass('d-block')
          }
        }
      })
    }
  })

  $('#count').on('change', function() {
    $('#countErr.invalid-feedback').removeClass('d-block')
  })

  $('#finishedProductName').on('keyup', function() {
    $('#noNameErr.invalid-feedback').removeClass('d-block')
  })

  $('#add').on('click', function() {
    let selectedItemName = $('#itemType option:selected')
    var name = selectedItemName.text()
    let itemCount = $('#count')
    var count = itemCount.val()
    if (!count) {
      $('#countErr.invalid-feedback').addClass('d-block')
      return
    }
    if (!name) {
      $('#nameErr.invalid-feedback').addClass('d-block')
      return
    }
    items.row
      .add({
        name: name,
        count: count
      })
      .draw()
    selectedItemName.remove()
    let remaining = $('#itemType').find('option').length
    if (remaining === 0) {
      $('#add').prop('disabled', true)
    }
    $('#noRowsErr.invalid-feedback').removeClass('d-block')
    itemCount.val('')
  })
})
