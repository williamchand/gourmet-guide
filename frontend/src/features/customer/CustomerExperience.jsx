import { useMemo, useState } from 'react';
import {
  DEFAULT_SESSION,
  RESTAURANTS,
  RESTAURANT_MENUS
} from './customerData.js';

function cloneSession() {
  return {
    ...DEFAULT_SESSION,
    allergies: [...DEFAULT_SESSION.allergies],
    transcript: DEFAULT_SESSION.transcript.map((entry) => ({ ...entry })),
    orderItems: []
  };
}

function formatPrice(value) {
  return `$${value.toFixed(2)}`;
}

export function CustomerExperience() {
  const [restaurantId, setRestaurantId] = useState(RESTAURANTS[0].id);
  const [isCheckoutOpen, setIsCheckoutOpen] = useState(false);
  const [selectedMenuTitle, setSelectedMenuTitle] = useState('');
  const [checkoutStatus, setCheckoutStatus] = useState('');
  const [sessionsByRestaurant, setSessionsByRestaurant] = useState(() =>
    RESTAURANTS.reduce(
      (accumulator, restaurant) => ({
        ...accumulator,
        [restaurant.id]: cloneSession()
      }),
      {}
    )
  );

  const currentSession = sessionsByRestaurant[restaurantId] ?? cloneSession();
  const { allergies, transcript, orderItems } = currentSession;

  const updateSession = (updater) => {
    setSessionsByRestaurant((current) => ({
      ...current,
      [restaurantId]: updater(current[restaurantId] ?? cloneSession())
    }));
  };

  const recommendations = useMemo(() => {
    const restaurantMenu = RESTAURANT_MENUS[restaurantId] ?? [];
    return restaurantMenu.map((entry) => {
      const hasConflict = allergies.some((allergy) =>
        entry.notes.toLowerCase().includes(allergy.toLowerCase())
      );
      return { ...entry, safeToOrder: !hasConflict };
    });
  }, [allergies, restaurantId]);

  const selectedMenu =
    recommendations.find((entry) => entry.title === selectedMenuTitle) ?? null;

  const orderTotal = useMemo(
    () => orderItems.reduce((sum, item) => sum + item.price, 0),
    [orderItems]
  );

  const pushTranscript = (message) => {
    updateSession((session) => ({
      ...session,
      transcript: [...session.transcript, message]
    }));
  };

  const addToOrder = (menuItem) => {
    updateSession((session) => ({
      ...session,
      orderItems: [
        ...session.orderItems,
        {
          title: menuItem.title,
          price: menuItem.price,
          dietaryTags: menuItem.dietaryTags
        }
      ],
      transcript: [
        ...session.transcript,
        {
          speaker: 'Assistant',
          text: `${menuItem.title} has been added to your order list.`
        }
      ]
    }));
  };

  const openDetails = (menuItem) => {
    setSelectedMenuTitle(menuItem.title);
    pushTranscript({
      speaker: 'Assistant',
      text: `Opened details for ${menuItem.title}.`
    });
  };

  const removeFromOrder = (indexToRemove) => {
    updateSession((session) => ({
      ...session,
      orderItems: session.orderItems.filter(
        (_, index) => index !== indexToRemove
      )
    }));
  };

  const finalizeOrder = () => {
    setCheckoutStatus('Your order is confirmed and sent to the kitchen.');
    setIsCheckoutOpen(false);
    pushTranscript({
      speaker: 'Assistant',
      text: 'Great choice. Your order has been finalized and sent to the kitchen.'
    });
    updateSession((session) => ({ ...session, orderItems: [] }));
  };

  return (
    <section className="customer-layout customer-layout--immersive">
      <article className="card customer-toolbar">
        <div>
          <h2>Customer Experience</h2>
          <p className="caption">Restaurant-specific voice ordering session.</p>
        </div>
        <div>
          <label htmlFor="restaurant-picker">Active Restaurant</label>
          <select
            id="restaurant-picker"
            value={restaurantId}
            onChange={(event) => {
              setRestaurantId(event.target.value);
              setCheckoutStatus('');
              setSelectedMenuTitle('');
            }}
          >
            {RESTAURANTS.map((restaurant) => (
              <option key={restaurant.id} value={restaurant.id}>
                {restaurant.name}
              </option>
            ))}
          </select>
        </div>
      </article>

      <article className="card customer-column customer-column--menu">
        <div className="recommendation-header">
          <h2>Voice-ranked Recommendations</h2>
          <p className="caption">Say: “Open details for Salmon Plate”.</p>
        </div>
        <div className="recommendation-list recommendation-list--tablet-fit">
          {recommendations.map((entry) => (
            <section
              key={entry.title}
              className="recommendation recommendation--tablet"
            >
              <img
                className="recommendation__image"
                src={entry.image}
                alt={entry.title}
              />
              <div className="recommendation__body">
                <h3>{entry.title}</h3>
                <div className="tag-row">
                  {entry.dietaryTags.map((tag) => (
                    <span key={tag} className="dietary-tag">
                      {tag}
                    </span>
                  ))}
                </div>
                <p className="price">{formatPrice(entry.price)}</p>
                <div className="recommendation__actions">
                  <button
                    type="button"
                    className="chip"
                    onClick={() => openDetails(entry)}
                  >
                    Open details
                  </button>
                  <button
                    type="button"
                    className="primary-button"
                    onClick={() => addToOrder(entry)}
                    disabled={!entry.safeToOrder}
                  >
                    {entry.safeToOrder ? 'Add to order' : 'Review allergy'}
                  </button>
                </div>
              </div>
            </section>
          ))}
        </div>
      </article>

      <article className="card customer-column customer-column--order">
        <h2>Order List</h2>
        {orderItems.length === 0 ? (
          <p className="caption">
            No items yet. Add dishes from recommendations.
          </p>
        ) : (
          <ul className="order-list">
            {orderItems.map((item, index) => (
              <li key={`${item.title}-${index}`} className="order-list__item">
                <div>
                  <p className="order-list__name">{item.title}</p>
                  <div className="tag-row">
                    {item.dietaryTags.map((tag) => (
                      <span
                        key={`${item.title}-${tag}-${index}`}
                        className="dietary-tag"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                  <p className="price">{formatPrice(item.price)}</p>
                </div>
                <button
                  type="button"
                  className="chip"
                  onClick={() => removeFromOrder(index)}
                >
                  Remove
                </button>
              </li>
            ))}
          </ul>
        )}
        <p className="risk">Total: {formatPrice(orderTotal)}</p>
        <button
          type="button"
          className="primary-button"
          onClick={() => setIsCheckoutOpen(true)}
          disabled={orderItems.length === 0}
        >
          Finalize order
        </button>
        {checkoutStatus && <p className="caption">{checkoutStatus}</p>}
      </article>

      <article className="card customer-column customer-column--voice">
        <div className="voice-panel-header">
          <button
            type="button"
            className="chip chip--active listening-pill"
            aria-live="polite"
          >
            Listening now…
          </button>
          <div className="tag-row tag-row--compact">
            {allergies.map((allergy) => (
              <span key={allergy} className="dietary-tag dietary-tag--danger">
                {allergy}
              </span>
            ))}
          </div>
        </div>

        <h3>Voice Chat Flow</h3>
        <ul className="timeline timeline--fixed">
          {transcript.map((entry, index) => (
            <li
              key={`${entry.speaker}-${index}`}
              className={`message message--${entry.speaker.toLowerCase()}`}
            >
              <strong>{entry.speaker}</strong>
              <p>{entry.text}</p>
            </li>
          ))}
        </ul>
      </article>

      {selectedMenu && (
        <div
          className="modal-overlay"
          role="dialog"
          aria-modal="true"
          aria-label="Menu detail popup"
        >
          <section className="modal-card">
            <h3>{selectedMenu.title}</h3>
            <p>{selectedMenu.notes}</p>
            <div className="tag-row">
              {selectedMenu.dietaryTags.map((tag) => (
                <span key={tag} className="dietary-tag">
                  {tag}
                </span>
              ))}
            </div>
            <p className="price">Price: {formatPrice(selectedMenu.price)}</p>
            <div className="modal-actions">
              <button
                type="button"
                className="primary-button primary-button--active"
                onClick={() => {
                  addToOrder(selectedMenu);
                  setSelectedMenuTitle('');
                }}
              >
                Add to order
              </button>
              <button
                type="button"
                className="primary-button"
                onClick={() => setSelectedMenuTitle('')}
              >
                Back to recommendations
              </button>
            </div>
          </section>
        </div>
      )}

      {isCheckoutOpen && (
        <div
          className="modal-overlay"
          role="dialog"
          aria-modal="true"
          aria-label="Finalize order confirmation"
        >
          <section className="modal-card">
            <h3>Finalize your order?</h3>
            <p>You can confirm now, or continue thinking about the menu.</p>
            <div className="modal-actions">
              <button
                type="button"
                className="primary-button primary-button--active"
                onClick={finalizeOrder}
              >
                Confirm order
              </button>
              <button
                type="button"
                className="primary-button"
                onClick={() => setIsCheckoutOpen(false)}
              >
                Keep exploring menu
              </button>
            </div>
          </section>
        </div>
      )}
    </section>
  );
}
