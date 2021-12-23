/**
 * @license
 * Copyright 2021 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

interface Target {
  toString(): string;
}

function sayHello(to: Target): void {
  console.log('Hello, ' + to + '!');
}

const world = {
  name: 'World',
  toString(): string {
    return this.name;
  },
};

sayHello(world);
