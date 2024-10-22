const defaultField = {
  name: "",
  selector: "",
  type: "text",
};

const CLASS_NAMES = {
  mainContainer: "flex items-center gap-x-1 mb-1",
  input:
    "disabled:cursor-not-allowed disabled:opacity-75 focus:outline-none border-0 form-input rounded-md placeholder-gray-400 text-sm px-2.5 py-1.5 shadow-sm bg-white text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-primary-500 max-w-72",
  select:
    "disabled:cursor-not-allowed disabled:opacity-75 focus:outline-none border-0 form-select rounded-md text-sm px-2.5 py-[7px] shadow-sm bg-white text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-primary-500 max-w-72",
  addBtn:
    "font-medium rounded-md text-sm px-2.5 py-1.5 shadow-sm ring-1 ring-gray-300 text-gray-700 bg-gray-50 hover:bg-gray-100 focus:ring-2 focus:ring-primary-500 flex items-center disabled:opacity-75 disabled:cursor-not-allowed",
};

const FieldItem = {
  nameInput(field) {
    const nameInput = document.createElement("input");
    nameInput.type = "text";
    nameInput.value = field.name;
    nameInput.placeholder = "Field name";
    nameInput.required = true;
    nameInput.name = "name";
    nameInput.className = CLASS_NAMES.input;

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
    selectorInput.className = CLASS_NAMES.input;

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
    typeSelect.className = CLASS_NAMES.select;

    const options = [
      { value: "text", text: "Text" },
      { value: "nested", text: "Nested" },
      { value: "attribute", text: "Attribute" },
    ];

    options.forEach(({ value, text }) => {
      const option = document.createElement("option");
      option.value = value;
      option.innerText = text;
      option.selected = value === field.type;
      typeSelect.appendChild(option);
    });

    typeSelect.addEventListener("change", () => {
      field.type = typeSelect.value;
    });

    return typeSelect;
  },

  attributeInput(field) {
    const attributeInput = document.createElement("input");
    attributeInput.value = field.attribute;
    attributeInput.placeholder = "Attribute name";
    attributeInput.required = true;
    attributeInput.name = "attribute";
    attributeInput.className = CLASS_NAMES.input;

    attributeInput.addEventListener("input", () => {
      field.attribute = attributeInput.value;
    });

    return attributeInput;
  },

  addFieldButton() {
    const addFieldButton = document.createElement("button");
    addFieldButton.innerText = "Add +";
    addFieldButton.className = CLASS_NAMES.addBtn;
    return addFieldButton;
  },

  childrenContainer() {
    const nestedFieldsContainer = document.createElement("div");
    nestedFieldsContainer.className = "ml-8 mb-1";
    return nestedFieldsContainer;
  },

  deleteFieldButton() {
    const deleteButton = document.createElement("button");
    deleteButton.innerHTML = `
<svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 10.5 10.5">
  <path d="M5.244 0.87A0.35 0.35 0 0 0 4.9 1.225V1.4H2.975a0.35 0.35 0 0 0 -0.355 0.35H2.1a0.35 0.35 0 1 0 0 0.7h6.3a0.35 0.35 0 1 0 0 -0.7h-0.52A0.35 0.35 0 0 0 7.525 1.4H5.6v-0.175a0.35 0.35 0 0 0 -0.356 -0.355M2.1 3.15l0.628 5.332A0.699 0.699 0 0 0 3.422 9.1h3.656a0.7 0.7 0 0 0 0.694 -0.618L8.4 3.15z"/>
</svg>
    `;

    deleteButton.className = "fill-red-500 ";

    return deleteButton;
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
  console.log("renderField_");

  // Clear the container
  container.innerHTML = "";
  const reRender = () => renderField(field, container, onDelete);

  // Create the first row: name, selector, type, attribute, delete btn (if type is attribute)
  const firstRow = document.createElement("div");
  firstRow.className = CLASS_NAMES.mainContainer;

  const nameInput = FieldItem.nameInput(field);
  const selectorInput = FieldItem.selectorInput(field);
  const typeSelect = FieldItem.typeSelect(field);
  const deleteButton = FieldItem.deleteFieldButton();

  firstRow.appendChild(nameInput);
  firstRow.appendChild(selectorInput);
  firstRow.appendChild(typeSelect);
  field.type === "attribute" &&
    firstRow.appendChild(FieldItem.attributeInput(field));
  firstRow.appendChild(deleteButton);

  // Append the first row to the container
  container.appendChild(firstRow);

  // Append the second row: nested fields
  if (field.type === "nested") {
    const childContainer = FieldItem.childrenContainer();
    container.appendChild(childContainer);
    renderFields(field.fields, childContainer);
  }

  // Add event listeners
  deleteButton.addEventListener("click", onDelete);
  typeSelect.addEventListener("change", () => {
    delete field.attribute;
    delete field.fields;

    if (field.type === "nested") {
      field.fields = [{ ...defaultField }];
    } else if (field.type === "attribute") {
      field.attribute = "";
    }

    reRender();
  });
}

function renderFields(fields, container) {
  container.innerHTML = "";
  const reRender = () => renderFields(fields, container);

  fields.forEach((field, i) => {
    const fieldContainer = document.createElement("div");
    container.appendChild(fieldContainer);

    renderField(field, fieldContainer, () => {
      fields.splice(i, 1);
      reRender();
    });
  });

  const addFieldButton = FieldItem.addFieldButton();
  addFieldButton.addEventListener("click", () => {
    fields.push({ ...defaultField });
    reRender();
  });
  container.appendChild(addFieldButton);
}
