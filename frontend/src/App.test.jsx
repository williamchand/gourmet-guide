import { fireEvent, render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import App from './App.jsx';

async function waitForVoiceConfigLoad() {
  await screen.findByText(/Gemini voice stream/i);
}

describe('App', () => {
  beforeEach(() => {
    window.history.pushState({}, '', '/');
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        model: 'gemini-2.5-flash-native-audio-preview-12-2025',
        config: {
          response_modalities: ['AUDIO'],
          system_instruction: 'You are a helpful and friendly AI assistant.'
        },
        audio: {
          send_sample_rate: 16000,
          receive_sample_rate: 24000
        }
      })
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('renders chooser by default and navigates to customer', async () => {
    render(<App />);

    expect(screen.getByText('Choose your workspace.')).toBeInTheDocument();
    fireEvent.click(
      screen.getByRole('button', { name: 'Open customer screen' })
    );

    await waitForVoiceConfigLoad();
    expect(window.location.pathname).toBe('/customer');
    expect(screen.getByText('Customer Experience')).toBeInTheDocument();
    expect(screen.getByText('Listening now…')).toBeInTheDocument();
    expect(
      screen.getByText('Voice-ranked Recommendations')
    ).toBeInTheDocument();
    expect(
      screen.getByText(/gemini-2.5-flash-native-audio-preview-12-2025/i)
    ).toBeInTheDocument();
  });

  it('navigates to admin from chooser', () => {
    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: 'Open admin screen' }));

    expect(window.location.pathname).toBe('/admin');
    expect(screen.getByText('Admin Access')).toBeInTheDocument();
    expect(screen.getByText('Combo Builder')).toBeInTheDocument();
  });

  it('loads admin experience from admin route', () => {
    window.history.pushState({}, '', '/admin');

    render(<App />);

    expect(screen.getByText('Menu Management Dashboard')).toBeInTheDocument();
  });

  it('keeps always-on voice flow and logs interactions', async () => {
    window.history.pushState({}, '', '/customer');
    render(<App />);

    await waitForVoiceConfigLoad();
    expect(screen.getByText('Peanuts')).toBeInTheDocument();
    fireEvent.click(screen.getAllByRole('button', { name: 'Add to order' })[0]);
    expect(
      screen.getByText(/has been added to your order list/)
    ).toBeInTheDocument();
  });

  it('supports restaurant-specific sessions and menus', async () => {
    window.history.pushState({}, '', '/customer');
    render(<App />);

    await waitForVoiceConfigLoad();
    expect(screen.getByText('Listening now…')).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('Active Restaurant'), {
      target: { value: 'green-garden' }
    });

    expect(screen.getByText('Tofu Lettuce Wraps')).toBeInTheDocument();
    expect(
      screen.queryByText('Herb-Roasted Salmon Plate')
    ).not.toBeInTheDocument();
  });

  it('opens menu detail popup and returns to recommendations', async () => {
    window.history.pushState({}, '', '/customer');
    render(<App />);

    await waitForVoiceConfigLoad();
    fireEvent.click(screen.getAllByRole('button', { name: 'Open details' })[0]);
    expect(
      screen.getByRole('dialog', { name: 'Menu detail popup' })
    ).toBeInTheDocument();

    fireEvent.click(
      screen.getByRole('button', { name: 'Back to recommendations' })
    );
    expect(
      screen.queryByRole('dialog', { name: 'Menu detail popup' })
    ).not.toBeInTheDocument();
  });

  it('supports order list and finalize order modal', async () => {
    window.history.pushState({}, '', '/customer');
    render(<App />);

    await waitForVoiceConfigLoad();
    fireEvent.click(screen.getAllByRole('button', { name: 'Add to order' })[0]);
    expect(screen.getByText('Total: $24.00')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Finalize order' }));
    expect(
      screen.getByRole('dialog', { name: 'Finalize order confirmation' })
    ).toBeInTheDocument();

    fireEvent.click(
      screen.getByRole('button', { name: 'Keep exploring menu' })
    );
    expect(
      screen.queryByRole('dialog', { name: 'Finalize order confirmation' })
    ).not.toBeInTheDocument();
  });

  it('runs the admin journey: login + setup + combo update', () => {
    window.history.pushState({}, '', '/admin');
    render(<App />);

    fireEvent.click(
      screen.getByRole('button', { name: 'Continue with Google' })
    );
    expect(screen.getByText(/JWT session active/)).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('Restaurant name'), {
      target: { value: 'Hackathon Bistro' }
    });
    fireEvent.change(screen.getByLabelText('Restaurant ID / slug'), {
      target: { value: 'hackathon-bistro' }
    });
    expect(screen.getByText(/Setup ready for/)).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('Combo name'), {
      target: { value: 'Hackathon Special Combo' }
    });
    expect(screen.getByText('Hackathon Special Combo')).toBeInTheDocument();
  });
});
