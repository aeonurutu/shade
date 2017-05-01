// Copyright 2016-2017 Richard Hawkins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Toolbar will be responsible for sending commands back to the app.
// For example, to pause, resume, ...
// Currently this is just a placeholder.
var Toolbar = React.createClass({
	name: "Toolbar",

	propTypes: function() {
		return {
			defaultEngineState: React.PropTypes.string,
		};
	},

	getDefaultProps: function() {
		return {
			defaultEngineState: null,
		};
	},

	getInitialState: function() {
		return {
			engineState: this.props.defaultEngineState,
		};
	},

	_play_pause_state: function(ev) {
		console.log("play or pause state");
	},

	_step_state: function(ev) {
		console.log("step state");
	},

	_download_state: function(ev) {
		// TODO: figure out why the state is only sometimes downloaded.
		var contents = JSON.stringify(this.state.engineState);
		console.log(contents);
		var URL = window.URL || window.webkit.URL;
		var blob = new Blob([contents], {type: 'text/json'});
		ev.target.href = URL.createObjectURL(blob);
		ev.target.download = this.state.engineState.project.name + '-state.json';
	},

	_upload_state: function(ev) {
		console.log("upload state");
	},

	componentWillReceiveProps: function(newProps) {
		this.setState({
			engineState: newProps.defaultEngineState,
		});
	},

	render: function() {
		return React.DOM.div({className: "btn-toolbar"}, 
			React.DOM.div({className: "btn-group"},
				// TODO: swap with glyphicon-play when paused.
				React.DOM.a({
					className: "btn btn-default",
					title: "pause/play game state",
					onClick: this._play_pause_state,
					href: "#",
				}, React.DOM.i({className: "glyphicon glyphicon-pause"}, "")),
				React.DOM.a({
					className: "btn btn-default",
					title: "pause game state immediatly following next Render call",
					onClick: this._step_state,
					href: "#",
				},React.DOM.i({className: "glyphicon glyphicon-step-forward"}, ""))
			),
			React.DOM.div({className: "btn-group"},
				React.DOM.a({
					className: "btn btn-default",
					title: "download state",
					onClick: this._download_state,
					href: "#"
				}, React.DOM.i({className: "glyphicon glyphicon-floppy-save"}, "")),
				React.DOM.a({
					className: "btn btn-default",
					title: "upload state, and replace in the engine",
					onClick: this._upload_state,
					href: "#"
				}, React.DOM.i({className: "glyphicon glyphicon-floppy-open"}, ""))
			)
		);
	},
});

// StateView is a simple way to dump all recieved engine state.  Eventually
// components should be made for all the individual pieces of state but since
// the format of this might change, just dumping it all as JSON is probably
// the best option right now.
var StateView = React.createClass({
	name: "StateView",

	propTypes: function() {
		return {
			defaultEngineState: React.PropTypes.string,
		};
	},

	getDefaultProps: function() {
		return {
			defaultEngineState: null,
		};
	},

	getInitialState: function() {
		return {
			engineState: this.props.defaultEngineState,
		};
	},

	componentWillReceiveProps: function(newProps) {
		this.setState({
			engineState: newProps.defaultEngineState,
		});
	},

	render: function() {
		return React.DOM.span(null,
			React.DOM.h1(null, "State"),
			React.DOM.pre({className: ".pre-scrollable"}, JSON.stringify(this.state.engineState, null, 2))
		);
	},
});

var DebugClient = React.createClass({
	name: "DebugClient",

	propTypes: function() {
		return {
			defaultEngineState: React.PropTypes.object,
		};
	},

	getDefaultProps: function () {
		return {
			defaultEngineState: {'foo': 'bar'},
		};
	},

	getInitialState: function() {
		return {
			socket: null,
			engineState: this.props.defaultEngineState,
		};
	},

	componentWillMount: function() {
		var comp = this;
		var ws_uri = "ws://" + window.location.host + "/debug";
		this.state.conn = new WebSocket(ws_uri, "protocolOne"); 

		this.state.conn.onopen = function(ev) {
	 		var msg = {
				method: "state",
			};
			this.send(JSON.stringify(msg));
		};

		this.state.conn.onclose = function(ev) {
			console.log(ev);
		};

		this.state.conn.onmessage = function(ev) {
			var edata = JSON.parse(ev.data);
			comp.setState({
				engineState: edata,
			});
		};

		this.state.conn.onerror = function(ev) {
		};
	},

	componentWillUnmount: function() {
		this.state.conn.close();
	},

	render: function() {
		return React.DOM.span(null,
			React.createElement(Toolbar, {defaultEngineState: this.state.engineState}),
			React.createElement(StateView, {defaultEngineState: this.state.engineState})
		);
	},
});

ReactDOM.render(
	React.createElement(DebugClient, {}),
	document.getElementById("app")
);
