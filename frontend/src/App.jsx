import { useMemo, useState } from 'react';

const AVAILABLE_ALLERGIES = [
  'Peanuts',
  'Tree Nuts',
  'Dairy',
  'Egg',
  'Gluten',
  'Soy',
  'Shellfish',
  'Sesame',
];

const MENU_ITEMS = [
  {
    id: 'item-1',
    name: 'Mediterranean Bowl',
    ingredients: 'Quinoa, hummus, cucumber, olives',
    allergens: ['Sesame'],
  },
  {
    id: 'item-2',
    name: 'Crispy Tofu Wrap',
    ingredients: 'Tofu, flour tortilla, slaw, aioli',
    allergens: ['Soy', 'Gluten', 'Egg'],
  },
];

const BASE_RECOMMENDATIONS = [
  {
    title: 'Herb-Roasted Salmon Plate',
    notes: 'No peanut ingredients detected. Kitchen should use separate utensils.',
    risk: 'Low risk',
  },
  {
    title: 'Citrus Quinoa Salad',
    notes: 'Ask for dressing on the side to avoid dairy contamination.',
    risk: 'Moderate risk',
  },
  {
    title: 'Smoky Veggie Combo',
    notes: 'Contains sesame in the default sauce. Request allergen-safe substitute.',
    risk: 'Needs adjustment',
  },
];

function AllergySelector({ selected, onToggle }) {
  return (
    <div className="chip-group" role="group" aria-label="Allergy preferences">
      {AVAILABLE_ALLERGIES.map((allergy) => {
        const active = selected.includes(allergy);
        return (
          <button
            key={allergy}
            type="button"
            className={`chip ${active ? 'chip--active' : ''}`}
            onClick={() => onToggle(allergy)}
          >
            {allergy}
          </button>
        );
      })}
    </div>
  );
}

function CustomerExperience() {
  const [isListening, setIsListening] = useState(false);
  const [allergies, setAllergies] = useState(['Peanuts']);
  const [transcript, setTranscript] = useState([
    { speaker: 'Assistant', text: 'Hi! I can help you find safe menu items. What should we avoid?' },
    { speaker: 'Customer', text: 'No peanuts, and I prefer dairy-free meals.' },
  ]);
  const [uploadedImageName, setUploadedImageName] = useState('');

  const recommendations = useMemo(
    () =>
      BASE_RECOMMENDATIONS.map((entry) => {
        const hasConflict = allergies.some((allergy) => entry.notes.toLowerCase().includes(allergy.toLowerCase()));
        return {
          ...entry,
          risk: hasConflict ? 'Review required' : entry.risk,
        };
      }),
    [allergies],
  );

  const toggleAllergy = (next) => {
    setAllergies((current) =>
      current.includes(next) ? current.filter((value) => value !== next) : [...current, next],
    );
    setTranscript((current) => [
      ...current,
      { speaker: 'Assistant', text: `Updated allergy profile to include ${next}. Re-running safety checks now.` },
    ]);
  };

  const toggleListening = () => {
    setIsListening((value) => !value);
    setTranscript((current) => [
      ...current,
      {
        speaker: 'Assistant',
        text: !isListening
          ? 'Listening… tell me what dish you are considering.'
          : 'Microphone paused. I can continue through text recommendations.',
      },
    ]);
  };

  const onImageChange = (event) => {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    setUploadedImageName(file.name);
    setTranscript((current) => [
      ...current,
      {
        speaker: 'Assistant',
        text: `Analyzing ${file.name} for allergen signals in ingredients and cross-contact warnings.`,
      },
    ]);
  };

  return (
    <section className="panel-grid">
      <article className="card">
        <h2>Realtime Concierge</h2>
        <p>Microphone-driven assistant with live transcript updates.</p>
        <button
          type="button"
          className={`primary-button ${isListening ? 'primary-button--active' : ''}`}
          onClick={toggleListening}
        >
          {isListening ? 'Stop microphone' : 'Start microphone'}
        </button>

        <h3>Allergy Profile</h3>
        <AllergySelector selected={allergies} onToggle={toggleAllergy} />

        <h3>Menu Vision Safety Check</h3>
        <label htmlFor="menu-image" className="upload-label">
          Upload/capture menu image
        </label>
        <input
          id="menu-image"
          type="file"
          accept="image/*"
          capture="environment"
          onChange={onImageChange}
        />
        {uploadedImageName && <p className="caption">Queued image: {uploadedImageName}</p>}
      </article>

      <article className="card">
        <h2>Transcript Stream</h2>
        <ul className="timeline">
          {transcript.map((entry, index) => (
            <li key={`${entry.speaker}-${index}`}>
              <strong>{entry.speaker}:</strong> {entry.text}
            </li>
          ))}
        </ul>
      </article>

      <article className="card">
        <h2>Recommendations</h2>
        <div className="recommendation-list">
          {recommendations.map((entry) => (
            <section key={entry.title} className="recommendation">
              <h3>{entry.title}</h3>
              <p>{entry.notes}</p>
              <p className="risk">Status: {entry.risk}</p>
            </section>
          ))}
        </div>
      </article>
    </section>
  );
}

function AdminExperience() {
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

export default function App() {
  const [view, setView] = useState('customer');

  return (
    <main className="container">
      <h1>GourmetGuide</h1>
      <p>Allergen-aware recommendations powered by Gemini on Google Cloud.</p>

      <div className="view-toggle" role="tablist" aria-label="Experience toggle">
        <button
          type="button"
          role="tab"
          aria-selected={view === 'customer'}
          className={`tab-button ${view === 'customer' ? 'tab-button--active' : ''}`}
          onClick={() => setView('customer')}
        >
          Customer Experience
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={view === 'admin'}
          className={`tab-button ${view === 'admin' ? 'tab-button--active' : ''}`}
          onClick={() => setView('admin')}
        >
          Admin Experience
        </button>
      </div>

      {view === 'customer' ? <CustomerExperience /> : <AdminExperience />}
    </main>
  );
}
