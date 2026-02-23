import { useEffect, useState } from 'react';
import { DEFAULT_ROUTE, ROUTES } from './routes.js';

function normalizeRoute(pathname) {
  if (
    pathname === ROUTES.home ||
    pathname === ROUTES.customer ||
    pathname === ROUTES.admin
  ) {
    return pathname;
  }

  return DEFAULT_ROUTE;
}

export function useRoute() {
  const [route, setRoute] = useState(() =>
    normalizeRoute(window.location.pathname)
  );

  useEffect(() => {
    if (route !== window.location.pathname) {
      window.history.replaceState({}, '', route);
    }
  }, [route]);

  useEffect(() => {
    const onPopState = () => setRoute(normalizeRoute(window.location.pathname));
    window.addEventListener('popstate', onPopState);
    return () => window.removeEventListener('popstate', onPopState);
  }, []);

  const navigate = (nextRoute) => {
    const normalized = normalizeRoute(nextRoute);
    window.history.pushState({}, '', normalized);
    setRoute(normalized);
  };

  return { route, navigate };
}
