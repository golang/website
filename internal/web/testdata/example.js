/**
 * @license
 * Copyright 2021 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
function sayHello(to) {
  console.log("Hello, " + to + "!");
}
const world = {
  name: "World",
  toString() {
    return this.name;
  }
};
sayHello(world);
