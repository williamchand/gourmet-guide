import { useState } from 'react';
import { AllergySelector } from '../shared/AllergySelector.jsx';
import { MENU_ITEMS } from './adminData.js';

export function AdminExperience() {
  const [menuItems, setMenuItems] = useState(MENU_ITEMS);
  const [selectedMenuItem, setSelectedMenuItem] = useState(MENU_ITEMS[0].id);
  const [comboItems, setComboItems] = useState(['Mediterranean Bowl']);
  const [comboName, setComboName] = useState('Lunch Balance Combo');

  const selectedItem = menuItems.find((item) => item.id === selectedMenuItem);

  const toggleTag = (allergy) => {
    setMenuItems((current) =>
      current.map((item) => {
        if (item.id !== selectedMenuItem) {
          return item;
        }

        const allergens = item.allergens.includes(allergy)
          ? item.allergens.filter((value) => value !== allergy)
          : [...item.allergens, allergy];
        return { ...item, allergens };
      }),
    );
  };

  const toggleComboItem = (itemName) => {
    setComboItems((current) =>
      current.includes(itemName) ? current.filter((name) => name !== itemName) : [...current, itemName],
    );
  };

  return (
    <section className="panel-grid">
      <article className="card">
        <h2>Menu Management Dashboard</h2>
        <p>Maintain canonical menu catalog with ingredient and allergen metadata.</p>
        <label htmlFor="menu-item-select">Menu item</label>
        <select
          id="menu-item-select"
          value={selectedMenuItem}
          onChange={(event) => setSelectedMenuItem(event.target.value)}
        >
          {menuItems.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
            </option>
          ))}
        </select>

        <h3>Ingredient Notes</h3>
        <textarea value={selectedItem?.ingredients ?? ''} readOnly rows={3} />

        <h3>Allergen Tags</h3>
        <AllergySelector selected={selectedItem?.allergens ?? []} onToggle={toggleTag} />
      </article>

      <article className="card">
        <h2>Combo Builder</h2>
        <label htmlFor="combo-name">Combo name</label>
        <input
          id="combo-name"
          type="text"
          value={comboName}
          onChange={(event) => setComboName(event.target.value)}
        />
        <p className="caption">Select items to include:</p>
        {menuItems.map((item) => (
          <label key={item.id} className="checkbox-row">
            <input
              type="checkbox"
              checked={comboItems.includes(item.name)}
              onChange={() => toggleComboItem(item.name)}
            />
            {item.name}
          </label>
        ))}
      </article>

      <article className="card">
        <h2>Combo Preview</h2>
        <p>
          <strong>{comboName}</strong>
        </p>
        <ul>
          {comboItems.map((item) => (
            <li key={item}>{item}</li>
          ))}
        </ul>
      </article>
    </section>
  );
}
