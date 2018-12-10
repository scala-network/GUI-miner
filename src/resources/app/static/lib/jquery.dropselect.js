/*!
* dropselect - v0.1
*
* jQuery plugin to make Bootstrap dropdown to "behave like" html select element
* Requires jQuery, Bootstrap
*
* Note: it doesn't add a hidden input. if you submit a form with Bootstrap dropdown,
* 		you will NOT have a value for the emulated dropdown
*
* Copyright (c) 2018 BLOC team.
* Inspired by Volkan Kucukcakar dropselect jQuery plugin
* https://github.com/vkucukcakar/dropselect
*/
$.fn.extend({
    dropselect: function(method, param) {
        switch (method) {

            // Constructor: Example: $('.dropselect').dropselect();
            case undefined:
            case '':
				return this.each(function() {
					var el = $(this);
					var button_text = el.find('span').first();
					var dropdown_menu = el.parent().find('.dropdown-menu');
					var dropdown_links = dropdown_menu.find('a');

					// make the dropdown behave like a normal select
					dropdown_links.off('click').on('click', function() {
						dropdown_links.removeClass('selected');
						$(this).addClass('selected');
						button_text.html($(this).html());
						el.trigger('dropdown-event');
					});

					// make the first value selected if no value is actually selected
					if (dropdown_menu.find('a.selected').length == 0) {
						dropdown_links.first().addClass('selected');
					}
				});

            // Manually select by value. Example: $('.dropselect').dropselect('select', 'anyval');
            case 'select':
				var dropdown_menu = $(this).parent().find('.dropdown-menu');
				dropdown_menu.find('a[data-value="' + param + '"]').trigger('click');
                return this;

            // Get selected option value. Example: var anyvalue = $('#mydropselect').dropselect('value');
            case 'value':
				var dropdown_menu = $(this).parent().find('.dropdown-menu');
				return dropdown_menu.find('a.selected').data('value');

            // Callback function when the select changes its value
            case 'change':
				if (typeof param == 'function') {
					this.on('dropdown-event', param);
				}
				return this;
        }
    }
});