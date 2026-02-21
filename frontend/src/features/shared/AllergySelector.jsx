import { AVAILABLE_ALLERGIES } from './allergyData.js';

export function AllergySelector({ selected, onToggle }) {
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
