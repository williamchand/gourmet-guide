import { render, screen } from '@testing-library/react';
import App from './App.jsx';

describe('App', () => {
  it('renders app headline', () => {
    render(<App />);
    expect(screen.getByText('GourmetGuide')).toBeInTheDocument();
  });
});
