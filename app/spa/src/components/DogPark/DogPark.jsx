import React, { useState } from "react";

import "./DogPark.scss";
import Dog from "./Dog";

const GRID_SIZE = 8;

export default function DogPark() {
  // const [grid, setGrid] = useState(Array(GRID_SIZE).fill(null).map(row => }{new Array(GRID_SIZE).fill(null)));
  const [grid, setGrid] = useState([
    ['husky', null, null, null, null, null, null, null],
    [null, null, null, null, null, null, null, null],
    [null, null, null, null, null, null, null, null],
    [null, null, null, null, 'boxer', null, null, null],
    [null, null, null, null, null, null, null, null],
    [null, null, null, null, null, null, null, null],
    [null, null, null, null, null, null, null, null],
    [null, null, null, null, null, null, null, 'chiuahua'],
  ]);

  return (
    <div>
      <div className="enclosure">
        <div className="tabs">
          <button className="tab active">
            <div>
              <span>Doggies Enclosure</span>
            </div>
          </button>
        </div>
        <div className="playground">
          <div className="playground__inner" id="tab-1">
            <div className="grid">
              {grid.map((row, x) => (
                row.map((elem, y) => elem ? <Dog key={`${x}${y}`} type={elem} position={{x, y}} /> : null)
              ))}
            </div>
          </div>
        </div>
      </div>
      <h1 className="mt-16">Hot doggos are ready to roam in the park</h1>
    </div>
  );
}
