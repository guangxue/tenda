import {
	ElementWrapper as $,
	Selection,
	WPButton,
} from './core.js';

function Sidebar(sideoption) {
	if(sideoption == 'left') {
		this.sidebar = $('.blog-left-menu');
	}
	if(sideoption == 'right') {
		this.sidebar = $('.blog-right-menu');
	}
};
Sidebar.prototype = WPButton;
Sidebar.prototype.addHeader = function(str) {
	this.sidebar
		.addChild( $('@div.panel-header').text(str).css("border-bottom:1px solid #ddd;") )
	return this;
};
Sidebar.prototype.addTitle = function(str) {
	this.sidebar
			.addChild( $('@div.panel-title').text(str) )
	return this;
};
Sidebar.prototype.addButtons = function(buttons) {
	this.sidebuttons = WPButton.build(buttons).wrapperElement('div').wrapperClass('sidebuttons');
	this.sidebar.addChild(this.sidebuttons.wrapper);
	return this;
};

export default Sidebar;