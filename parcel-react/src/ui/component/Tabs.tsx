import * as React from 'react'

export const Tabs = () => {
  return (
      <div className="py-1">
        <ul className="nav nav-tabs">
          <li className="nav-item">
            <a className="nav-link active disabled" href="#">Inventory</a>
          </li>
          <li className="nav-item">
            <a className="nav-link" href="config.html">Finished products</a>
          </li>
          <li className="nav-item">
            <a className="nav-link" href="stock.html">Stock</a>
          </li>
        </ul>
      </div>
  )
}