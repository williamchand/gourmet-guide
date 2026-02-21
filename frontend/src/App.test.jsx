import { fireEvent, render, screen } from '@testing-library/react';
import App from './App.jsx';

describe('App', () => {
  beforeEach(() => {
    window.history.pushState({}, '', '/customer');
  });

  it('renders customer experience by default', () => {
    render(<App />);

    expect(screen.getByText('Realtime Concierge')).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'Start microphone' })
    ).toBeInTheDocument();
    expect(screen.getByText('Transcript Stream')).toBeInTheDocument();
  });

  it('switches to admin dashboard experience and updates route', () => {
    render(<App />);

    fireEvent.click(screen.getByRole('tab', { name: 'Admin Experience' }));

    expect(window.location.pathname).toBe('/admin');
    expect(screen.getByText('Menu Management Dashboard')).toBeInTheDocument();
    expect(screen.getByText('Combo Builder')).toBeInTheDocument();
  });

  it('loads admin experience from admin route', () => {
    window.history.pushState({}, '', '/admin');

    render(<App />);

    expect(screen.getByText('Menu Management Dashboard')).toBeInTheDocument();
  });

  it('runs the customer journey: microphone + allergy updates + vision upload', () => {
    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: 'Start microphone' }));
    expect(
      screen.getByText(/Listening… tell me what dish/)
    ).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Shellfish' }));
    expect(
      screen.getByText(/Updated allergy profile to include Shellfish/)
    ).toBeInTheDocument();

    const file = new File(['menu image bytes'], 'menu-photo.png', {
      type: 'image/png'
    });
    fireEvent.change(screen.getByLabelText('Upload/capture menu image'), {
      target: { files: [file] }
    });

    expect(screen.getByText('Queued image: menu-photo.png')).toBeInTheDocument();
    expect(
      screen.getByText(/Analyzing menu-photo.png for allergen signals/)
    ).toBeInTheDocument();
  });

  it('runs the admin journey: update combo and preview', () => {
    window.history.pushState({}, '', '/admin');
    render(<App />);

    fireEvent.change(screen.getByLabelText('Combo name'), {
      target: { value: 'Hackathon Special Combo' }
    });
    expect(screen.getByText('Hackathon Special Combo')).toBeInTheDocument();

    const itemToggle = screen.getByRole('checkbox', {
      name: 'Citrus Herb Chicken'
    });
    fireEvent.click(itemToggle);

    expect(screen.getByRole('listitem', { name: 'Citrus Herb Chicken' })).toBeInTheDocument();
  });
});
