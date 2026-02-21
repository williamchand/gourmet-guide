import { AdminExperience } from '../features/admin/AdminExperience.jsx';
import { CustomerExperience } from '../features/customer/CustomerExperience.jsx';
import { ROUTES } from './routes.js';
import { useRoute } from './useRoute.js';

function RouteTabs({ route, onNavigate }) {
  return (
    <div className="view-toggle" role="tablist" aria-label="Experience toggle">
      <button
        type="button"
        role="tab"
        aria-selected={route === ROUTES.customer}
        className={`tab-button ${route === ROUTES.customer ? 'tab-button--active' : ''}`}
        onClick={() => onNavigate(ROUTES.customer)}
      >
        Customer Experience
      </button>
      <button
        type="button"
        role="tab"
        aria-selected={route === ROUTES.admin}
        className={`tab-button ${route === ROUTES.admin ? 'tab-button--active' : ''}`}
        onClick={() => onNavigate(ROUTES.admin)}
      >
        Admin Experience
      </button>
    </div>
  );
}

export function App() {
  const { route, navigate } = useRoute();

  return (
    <main className="container">
      <h1>GourmetGuide</h1>
      <p>Allergen-aware recommendations powered by Gemini on Google Cloud.</p>

      <RouteTabs route={route} onNavigate={navigate} />
      {route === ROUTES.admin ? <AdminExperience /> : <CustomerExperience />}
    </main>
  );
}
