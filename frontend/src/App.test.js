import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

describe('App Component', () => {
  test('renders the correct heading', () => {
    render(<App />);
    const headingElement = screen.getByText(/Inteli - Instituto de Tecnologia e Liderança/i);
    expect(headingElement).toBeInTheDocument();
  });

  test('renders the correct paragraph', () => {
    render(<App />);
    const paragraphElement = screen.getByText(/Já se fazem faculdades como futuramente\./i);
    expect(paragraphElement).toBeInTheDocument();
  });
});
