import { fireEvent, render, screen } from '@testing-library/react';
import App from './App.jsx';

describe('App', () => {
  it('renders customer experience by default', () => {
    render(<App />);

    expect(screen.getByText('Realtime Concierge')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Start microphone' })).toBeInTheDocument();
    expect(screen.getByText('Transcript Stream')).toBeInTheDocument();
  });

  it('switches to admin dashboard experience', () => {
    render(<App />);

    fireEvent.click(screen.getByRole('tab', { name: 'Admin Experience' }));

    expect(screen.getByText('Menu Management Dashboard')).toBeInTheDocument();
    expect(screen.getByText('Combo Builder')).toBeInTheDocument();
  });

  it('updates transcript when microphone is toggled', () => {
    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: 'Start microphone' }));

    expect(screen.getByText(/Listening… tell me what dish/)).toBeInTheDocument();
  });
});
