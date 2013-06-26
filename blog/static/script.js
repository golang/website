// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

$(function() {
	// Insert line numbers for all playground elements.
	$('.playground').each(function() {
		var $spans = $(this).find('> pre > span');

		// Compute width of number column (including trailing space).
		var max = 0;
		$spans.each(function() {
			var n = $(this).attr('num')*1;
			if (n > max) max = n;
		});
		var width = 2;
		while (max > 10) {
			max = max / 10;
			width++;
		}

		// Insert line numbers with space padding.
		$spans.each(function() {
			var n = $(this).attr('num')+" ";
			while (n.length < width) n = " "+n;
			$('<span class="number">').text(n).insertBefore(this);
		});
	});
});
