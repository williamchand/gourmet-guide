export const RESTAURANTS = [
  {
    id: 'harbor-grill',
    name: 'Harbor Grill'
  },
  {
    id: 'green-garden',
    name: 'Green Garden Cafe'
  }
];

const DEFAULT_TRANSCRIPT = [
  {
    speaker: 'Assistant',
    text: 'Hi! I can help you find safe menu items. What should we avoid?'
  },
  { speaker: 'Customer', text: 'No peanuts, and I prefer dairy-free meals.' },
  {
    speaker: 'Assistant',
    text: 'Listening… tell me what dish you are considering.'
  }
];

export const DEFAULT_SESSION = {
  isListening: true,
  allergies: ['Peanuts'],
  transcript: DEFAULT_TRANSCRIPT,
  orderItems: []
};

export const RESTAURANT_MENUS = {
  'harbor-grill': [
    {
      title: 'Herb-Roasted Salmon Plate',
      notes:
        'No peanut ingredients detected. Kitchen should use separate utensils and grill area.',
      risk: 'Low risk',
      dietaryTags: ['Gluten-free', 'High protein', 'Omega-3 rich'],
      image:
        'https://images.unsplash.com/photo-1467003909585-2f8a72700288?auto=format&fit=crop&w=900&q=80',
      price: 24
    },
    {
      title: 'Citrus Quinoa Salad',
      notes: 'Ask for dressing on the side to avoid dairy contamination.',
      risk: 'Moderate risk',
      dietaryTags: ['Vegan option', 'Dairy-aware', 'Fiber rich'],
      image:
        'https://images.unsplash.com/photo-1512621776951-a57141f2eefd?auto=format&fit=crop&w=900&q=80',
      price: 16
    },
    {
      title: 'Smoky Veggie Combo',
      notes:
        'Contains sesame in the default sauce. Request allergen-safe substitute.',
      risk: 'Needs adjustment',
      dietaryTags: ['Vegetarian', 'Customizable', 'Sesame alert'],
      image:
        'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?auto=format&fit=crop&w=900&q=80',
      price: 18
    }
  ],
  'green-garden': [
    {
      title: 'Tofu Lettuce Wraps',
      notes:
        'Prepared dairy-free by default. Confirm soy-only wok for severe allergy.',
      risk: 'Low risk',
      dietaryTags: ['Vegan', 'Dairy-free', 'High protein'],
      image:
        'https://images.unsplash.com/photo-1498837167922-ddd27525d352?auto=format&fit=crop&w=900&q=80',
      price: 14
    },
    {
      title: 'Miso Udon Bowl',
      notes:
        'Broth may include shellfish stock; request vegetable broth variant.',
      risk: 'Needs adjustment',
      dietaryTags: ['Vegetarian option', 'Warm bowl', 'Shellfish alert'],
      image:
        'https://images.unsplash.com/photo-1617093727343-374698b1b08d?auto=format&fit=crop&w=900&q=80',
      price: 17
    },
    {
      title: 'Coconut Mango Chia Cup',
      notes:
        'Naturally peanut-free and dairy-free. Built as a safe dessert recommendation.',
      risk: 'Low risk',
      dietaryTags: ['Dairy-free', 'Nut-aware', 'Dessert'],
      image:
        'https://images.unsplash.com/photo-1488477181946-6428a0291777?auto=format&fit=crop&w=900&q=80',
      price: 9
    }
  ]
};
