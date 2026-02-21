import { useMemo, useState } from 'react';
import { AllergySelector } from '../shared/AllergySelector.jsx';
import { BASE_RECOMMENDATIONS } from './customerData.js';

export function CustomerExperience() {
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
        return { ...entry, risk: hasConflict ? 'Review required' : entry.risk };
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
