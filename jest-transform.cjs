/**
 * @license
 * Copyright 2022 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

const {transform} = require('esbuild');

exports.createTransformer = () => ({
  canInstrument: true,
  processAsync: async (source) => {
    const result = await transform(source, {
      loader: 'ts',
    });
    if (result.warnings.length) {
      result.warnings.forEach(m => {
        console.warn(m);
      });
    }
    return {
      code: result.code,
      map: result.map,
    };
  },
});
