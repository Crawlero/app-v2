const defaultField = {
  name: "",
  selector: "",
  type: "text",
};

const FieldItem = {
  nameInput(field) {
    const nameInput = document.createElement("input", { type: "text" });
    nameInput.value = field.name;
    nameInput.placeholder = "Field name";
    nameInput.required = true;
    nameInput.name = "name";
    nameInput.className =
      "mt-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm";

    nameInput.addEventListener("input", () => {
      field.name = nameInput.value;
    });

    return nameInput;
  },
  selectorInput(field) {
    const selectorInput = document.createElement("input");
    selectorInput.value = field.selector;
    selectorInput.placeholder = "CSS selector";
    selectorInput.required = true;
    selectorInput.name = "selector";
    selectorInput.className =
      "mt-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm";

    selectorInput.addEventListener("input", () => {
      field.selector = selectorInput.value;
    });

    return selectorInput;
  },
  typeSelect(field) {
    const typeSelect = document.createElement("select");
    typeSelect.value = field.type;
    typeSelect.readOnly = true;
    typeSelect.name = "type";

    const textOption = document.createElement("option");
    textOption.value = "text";
    textOption.innerText = "Text";

    const nestedOption = document.createElement("option");
    nestedOption.value = "nested";
    nestedOption.innerText = "Nested";

    const attributeOption = document.createElement("option");
    attributeOption.value = "attribute";
    attributeOption.innerText = "Attribute";

    const options = [textOption, nestedOption, attributeOption];

    options.forEach((option) => {
      if (option.value === field.type) {
        option.selected = true;
      }

      typeSelect.appendChild(option);
    });

    typeSelect.className =
      "mt-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm";

    return typeSelect;
  },
  attributeInput(field) {
    const attributeInput = document.createElement("input");
    attributeInput.value = field.attribute;
    attributeInput.placeholder = "Attribute name";
    attributeInput.required = true;
    attributeInput.name = "attribute";
    attributeInput.className =
      "mt-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm";

    attributeInput.addEventListener("input", () => {
      field.attribute = attributeInput.value;
    });

    return attributeInput;
  },
  addFieldButton() {
    const addFieldButton = document.createElement("button");
    addFieldButton.innerText = "Add Field";
    addFieldButton.className =
      "mt-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm";

    return addFieldButton;
  },
  childrenContainer() {
    const nestedFieldsContainer = document.createElement("div");
    nestedFieldsContainer.className = "ml-8";
    return nestedFieldsContainer;
  },
  deleteFieldButton() {
    const deleteFieldButton = document.createElement("button");
    deleteFieldButton.innerText = "Delete Field";
    deleteFieldButton.className =
      "mt-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm";

    return deleteFieldButton;
  },
};

/**
 * @typedef {Object} TextField
 * @property {string} name - The name of the field.
 * @property {string} selector - The CSS selector to locate the element.
 * @property {"text"} type - The type of the field, which is "text".
 */

/**
 * @typedef {Object} NestedField
 * @property {string} name - The name of the field.
 * @property {string} selector - The CSS selector to locate the element.
 * @property {"nested"} type - The type of the field, which is "nested".
 * @property {Field[]} fields - Nested fields.
 */

/**
 * @typedef {Object} AttributeField
 * @property {string} name - The name of the field.
 * @property {string} selector - The CSS selector to locate the element.
 * @property {"attribute"} type - The type of the field, which is "attribute".
 * @property {string} attribute - The attribute to extract.
 */

/**
 * @typedef {TextField | NestedField | AttributeField} Field
 */

/**
 * Renders a field based on the field data.
 * @param {Field} field - The field data.
 * @param {HTMLDivElement} container - The container element to append the field.
 * @returns {HTMLDivElement} The field container.
 */
function renderField(field, container, onDelete) {
  container.innerHTML = "";
  const reRender = () => renderField(field, container);

  const mainContainer = document.createElement("div");
  const deleteButton = FieldItem.deleteFieldButton();
  const nameInput = FieldItem.nameInput(field);
  const selectorInput = FieldItem.selectorInput(field);
  const typeSelect = FieldItem.typeSelect(field, reRender);

  // append elements to the main container
  mainContainer.appendChild(deleteButton);
  mainContainer.appendChild(nameInput);
  mainContainer.appendChild(selectorInput);
  mainContainer.appendChild(typeSelect);
  if (field.type === "attribute")
    mainContainer.appendChild(FieldItem.attributeInput(field));

  // append the main container to the container
  container.appendChild(mainContainer);

  // add event listeners
  deleteButton.addEventListener("click", onDelete);
  typeSelect.addEventListener("change", () => {
    field.type = typeSelect.value;
    // reset
    delete field.attribute;
    delete field.fields;

    // set default values
    if (field.type === "nested") {
      field.fields = [defaultField];
    } else if (field.type === "attribute") {
      field.attribute = "";
    }
    reRender();
  });

  // render nested fields
  if (field.type === "nested") {
    const childContainer = FieldItem.childrenContainer();
    container.appendChild(childContainer);
    renderFields(field.fields, childContainer);
  }
}

/**
 * Renders a list of fields.
 * @param {Field[]} fields - The list of fields.
 * @param {HTMLDivElement} container - The container element to append the fields.
 */
function renderFields(fields, container) {
  container.innerHTML = "";

  const reRender = () => renderFields(fields, container);

  fields.forEach((field, i) => {
    const fieldContainer = document.createElement("div");
    renderField(field, fieldContainer, () => {
      // handle delete field
      fields.splice(i, 1);
      reRender();
    });
    container.appendChild(fieldContainer);
  });

  const addFieldButton = FieldItem.addFieldButton();
  addFieldButton.addEventListener("click", () => {
    fields.push({ ...defaultField });
    reRender();
  });
  container.appendChild(addFieldButton);
}
