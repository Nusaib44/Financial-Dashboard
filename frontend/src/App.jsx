import React, { useState, useEffect } from 'react';
import './App.css';

// Mock JWT for local development (skipping signature validation in backend, so structure matters)
// We need a structure that has "sub" and "email"
const MOCK_JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJlbWFpbCI6ImZvdW5kZXJAbGFnLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.NO_SIG";

function App() {
  const [agency, setAgency] = useState(null);
  const [loading, setLoading] = useState(true);
  const [showCreateForm, setShowCreateForm] = useState(false);

  // Form State
  const [name, setName] = useState('');
  const [currency, setCurrency] = useState('USD');
  const [startingCash, setStartingCash] = useState(0);

  useEffect(() => {
    fetchAgency();
  }, []);

  const fetchAgency = async () => {
    setLoading(true);
    try {
      const res = await fetch('/api/agency', {
        headers: {
          'CF-Access-Jwt-Assertion': MOCK_JWT // Injecting mock auth for now
        }
      });

      if (res.status === 200) {
        const data = await res.json();
        setAgency(data);
        setShowCreateForm(false);
      } else if (res.status === 404) {
        setAgency(null);
        setShowCreateForm(true);
      } else {
        console.error("Failed to fetch agency", res.status);
      }
    } catch (err) {
      console.error("Error fetching agency", err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateAgency = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch('/api/agency', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'CF-Access-Jwt-Assertion': MOCK_JWT
        },
        body: JSON.stringify({
          name,
          base_currency: currency,
          starting_cash: Number(startingCash)
        })
      });

      if (res.status === 201) {
        fetchAgency();
      } else if (res.status === 409) {
        alert("Agency already exists!");
      } else {
        alert("Failed to create agency");
      }
    } catch (err) {
      console.error("Error creating agency", err);
    }
  };

  if (loading) return <div>Loading...</div>;

  if (showCreateForm) {
    return (
      <div className="container">
        <h1>Create Your Agency</h1>
        <form onSubmit={handleCreateAgency}>
          <div>
            <label>Agency Name:</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          <div>
            <label>Base Currency:</label>
            <input
              type="text"
              value={currency}
              onChange={(e) => setCurrency(e.target.value)}
              required
            />
          </div>
          <div>
            <label>Starting Cash:</label>
            <input
              type="number"
              value={startingCash}
              onChange={(e) => setStartingCash(e.target.value)}
              required
            />
          </div>
          <button type="submit">Launch Agency</button>
        </form>
      </div>
    );
  }

  if (agency) {
    return (
      <div className="container">
        <h1>Agency Finance Reality — Founder Mode</h1>
        <h2>Agency: {agency.name}</h2>
      </div>
    );
  }

  return <div>Something went wrong.</div>;
}

export default App;
