import { AdminExperience } from '../features/admin/AdminExperience.jsx';
import { CustomerExperience } from '../features/customer/CustomerExperience.jsx';
import { ROUTES } from './routes.js';
import { useRoute } from './useRoute.js';

function ExperienceChooser({ onNavigate }) {
  return (
    <section className="experience-chooser">
      <h1>GourmetGuide</h1>
      <p>Choose your workspace.</p>
      <div className="experience-chooser__grid">
        <article className="card">
          <h2>Customer Experience</h2>
          <p>
            Voice-first tablet ordering with recommendations and live order
            summary.
          </p>
          <button
            type="button"
            className="primary-button primary-button--active"
            onClick={() => onNavigate(ROUTES.customer)}
          >
            Open customer screen
          </button>
        </article>
        <article className="card">
          <h2>Admin Experience</h2>
          <p>Google JWT sign-in simulation and restaurant setup dashboard.</p>
          <button
            type="button"
            className="primary-button"
            onClick={() => onNavigate(ROUTES.admin)}
          >
            Open admin screen
          </button>
        </article>
      </div>
    </section>
  );
}

function ExperienceHeader({ route, onNavigate }) {
  return (
    <header className="experience-header">
      <button
        type="button"
        className="chip"
        onClick={() => onNavigate(ROUTES.home)}
      >
        Back to chooser
      </button>
      <span className="caption">
        {route === ROUTES.customer ? 'Customer view' : 'Admin view'}
      </span>
    </header>
  );
}

export function App() {
  const { route, navigate } = useRoute();

  if (route === ROUTES.home) {
    return (
      <main className="container">
        <ExperienceChooser onNavigate={navigate} />
      </main>
    );
  }

  return (
    <main className="container container--immersive">
      <ExperienceHeader route={route} onNavigate={navigate} />
      {route === ROUTES.admin ? <AdminExperience /> : <CustomerExperience />}
    </main>
  );
}
