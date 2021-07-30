var JsonView = (function (exports) {
  'use strict';

  function _typeof(obj) {
    "@babel/helpers - typeof";

    if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") {
      _typeof = function (obj) {
        return typeof obj;
      };
    } else {
      _typeof = function (obj) {
        return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj;
      };
    }

    return _typeof(obj);
  }

  function expandedTemplate() {
    var params = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
    var key = params.key,
        size = params.size;
    return "\n    <div class=\"line\">\n      <div class=\"caret-icon\"><i class=\"fas fa-caret-right\"></i></div>\n      <div class=\"json-key json-expandable\">".concat(key, "</div>\n      <div class=\"json-size\">").concat(size, "</div>\n    </div>\n  ");
  }

  function notExpandedTemplate() {
    var params = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
    var key = params.key,
        value = params.value,
        type = params.type;
    var renderedValue = value;
    if (type === "object" && value && Object.keys(value).length < 1) {
      renderedValue = "{}";
    }
    if (type == "string") {
      renderedValue = renderedValue
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
      renderedValue = `"${renderedValue}"`;
    }
    return "\n    <div class=\"line\">\n      <div class=\"empty-icon\"></div>\n      <div class=\"json-key\">".concat(key, "</div>\n      <div class=\"json-separator\">:</div>\n      <div class=\"json-value json-").concat(type, "\">").concat(renderedValue, "</div>\n    </div>\n  ");
  }

  function hideNodeChildren(node) {
    node.children.forEach(function (child) {
      child.el.classList.add('hide');

      if (child.isExpanded) {
        hideNodeChildren(child);
      }
    });
  }

  function showNodeChildren(node) {
    node.children.forEach(function (child) {
      child.el.classList.remove('hide');

      if (child.isExpanded) {
        showNodeChildren(child);
      }
    });
  }

  function setCaretIconDown(node) {
    if (node.children.length > 0) {
      var icon = node.el.querySelector('.fas');

      if (icon) {
        icon.classList.replace('fa-caret-right', 'fa-caret-down');
      }
    }
  }

  function setCaretIconRight(node) {
    if (node.children.length > 0) {
      var icon = node.el.querySelector('.fas');

      if (icon) {
        icon.classList.replace('fa-caret-down', 'fa-caret-right');
      }
    }
  }

  function toggleNode(node, toggleAll) {
    if (node.isExpanded) {
      node.isExpanded = false;
      setCaretIconRight(node);
      hideNodeChildren(node);
      if (toggleAll) {
        collapseChildren(node);
      }
    } else {
      node.isExpanded = true;
      setCaretIconDown(node);
      showNodeChildren(node);
      if (toggleAll) {
        expandChildren(node);
      }
    }
  }

  function createContainerElement() {
    var el = document.createElement('div');
    el.className = 'json-container';
    return el;
  }

  function createNodeElement(node) {
    var el = document.createElement('div');

    var getSizeString = function getSizeString(node) {
      var len = node.children.length;
      if (node.type === 'array') return "[".concat(len, "]");
      if (node.type === 'object') return "{".concat(len, "}");
      return null;
    };

    if (node.children.length > 0) {
      el.innerHTML = expandedTemplate({
        key: node.key,
        size: getSizeString(node)
      });
      var caretEl = el.querySelector('.caret-icon');
      var expandedEl = el.querySelector('.json-key');
      [caretEl, expandedEl].forEach((cel) => cel.addEventListener('click', function (evt) {
        var toggleAll = evt.shiftKey;
        toggleNode(node, toggleAll);
      }));
    } else {
      el.innerHTML = notExpandedTemplate({
        key: node.key,
        value: node.value,
        type: _typeof(node.value)
      });
    }

    var lineEl = el.children[0];

    if (node.parent !== null) {
      lineEl.classList.add('hide');
    }

    lineEl.style = 'margin-left: ' + node.depth * 18 + 'px;';
    return lineEl;
  }

  function getDataType(val) {
    var type = _typeof(val);

    if (Array.isArray(val)) type = 'array';
    if (val === null) type = 'null';
    return type;
  }

  function traverseTree(node, callback) {
    callback(node);

    if (node.children.length > 0) {
      node.children.forEach(function (child) {
        traverseTree(child, callback);
      });
    }
  }

  function createNode() {
    var opt = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
    return {
      key: opt.key || null,
      parent: opt.parent || null,
      value: opt.hasOwnProperty('value') ? opt.value : null,
      isExpanded: opt.isExpanded || false,
      type: opt.type || null,
      children: opt.children || [],
      el: opt.el || null,
      depth: opt.depth || 0
    };
  }

  function createSubnode(data, node) {
    if (_typeof(data) === 'object') {
      for (var key in data) {
        var child = createNode({
          value: data[key],
          key: key,
          depth: node.depth + 1,
          type: getDataType(data[key]),
          parent: node
        });
        node.children.push(child);
        createSubnode(data[key], child);
      }
    }
  }

  function createTree(jsonData) {
    var data = parseTreeJSON(jsonData);
    var rootNode = createNode({
      value: data,
      key: getDataType(data),
      type: getDataType(data)
    });
    createSubnode(data, rootNode);
    return rootNode;
  }

  function renderJSON(jsonData, targetElement) {
    var tree = createTree(parseTreeJSON(jsonData));
    render(tree, targetElement);
    return tree;
  }

  function parseTreeJSON(jsonData) {
    if (jsonData && (typeof jsonData !== "string")) {
      return jsonData;
    }
    var parseData;
    try {
      parseData = JSON.parse(jsonData);
    } catch(e) {
      parseData = {
        error: `Could not parse JSON from data: ${e.message}`,
        data: jsonData
      };
    }
    return parseData;
  }

  function render(tree, targetElement) {
    var containerEl = createContainerElement();
    traverseTree(tree, function (node) {
      node.el = createNodeElement(node);
      containerEl.appendChild(node.el);
    });
    targetElement.appendChild(containerEl);
  }

  function expandChildrenDepth(node, depth) {
    node.el.classList.remove('hide');
    if (depth > 0 && node.children.length > 0) {
      node.isExpanded = true;
      setCaretIconDown(node);
      node.children.forEach(function (child) {
        expandChildrenDepth(child, depth - 1);
      });
    }
  }

  function expandChildren(node, depthLimit) {
    traverseTree(node, function (child) {
      child.el.classList.remove('hide');
      child.isExpanded = true;
      setCaretIconDown(child);
    });
  }

  function collapseChildren(node) {
    traverseTree(node, function (child) {
      child.isExpanded = false;
      if (child.depth > node.depth) child.el.classList.add('hide');
      setCaretIconRight(child);
    });
  }

  exports.collapseChildren = collapseChildren;
  exports.createTree = createTree;
  exports.expandChildren = expandChildren;
  exports.render = render;
  exports.renderJSON = renderJSON;
  exports.traverseTree = traverseTree;
  exports.expandChildrenDepth = expandChildrenDepth;

  return exports;

}({}));
