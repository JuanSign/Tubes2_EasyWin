'use client';

import dynamic from 'next/dynamic';
import { useState } from 'react';

// Dynamically imported components
const SelectInput = dynamic(() => import('../components/SelectInput'), { ssr: false });
const FlowDiagram = dynamic(() => import('../components/FlowDiagram'), { ssr: false });

export default function HomePage() {
  const [selected, setSelected] = useState('');
  const [showPopup, setShowPopup] = useState(false);
  const [algorithm, setAlgorithm] = useState('');
  const [mode, setMode] = useState('');
  const [diagramData, setDiagramData] = useState<any | null>(null);

  const handleFirstSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (selected) {
      setShowPopup(true);
    } else {
      alert('Please select an option first.');
    }
  };

  const handleFinalSearch = async () => {
    if (algorithm && mode) {
      const baseUrl = 'https://tubes2-easywin.onrender.com';
      const endpoint = algorithm.toLowerCase();
      const url = `${baseUrl}/${endpoint}`;

      try {
        const response = await fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            element: selected,
            type: mode === 'One Solution' ? 'one' : 'all',
          }),
        });

        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

        const data = await response.json();
        console.log('Response data:', data);
        setDiagramData(data);
        setShowPopup(false);
      } catch (err) {
        console.error('Fetch error:', err);
        alert('Failed to fetch data. Please try again.');
      }
    } else {
      alert('Please select both algorithm and mode.');
    }
  };

  const clearState = () => {
    setSelected('');
    setShowPopup(false);
    setAlgorithm('');
    setMode('');
    setDiagramData(null);
  };

  return (
    <main
      style={{
        display: 'flex',
        flexDirection: 'column',
        height: '100vh',
        backgroundColor: '#121212',
        color: '#f0f0f0',
        fontFamily: 'Segoe UI, sans-serif',
        padding: '20px',
      }}
    >
      <form
        onSubmit={handleFirstSubmit}
        style={{
          display: 'flex',
          gap: '10px',
          justifyContent: 'center',
          marginBottom: '20px',
        }}
      >
        <SelectInput onSelect={setSelected} />
        <button
          type="submit"
          style={{
            background: '#333',
            color: '#fff',
            padding: '8px 12px',
            border: 'none',
            borderRadius: 4,
            cursor: 'pointer',
          }}
        >
          Search
        </button>
        <button
          type="button"
          onClick={clearState}
          style={{
            background: '#555',
            color: '#fff',
            padding: '8px 12px',
            border: 'none',
            borderRadius: 4,
            cursor: 'pointer',
          }}
        >
          Clear
        </button>
      </form>

      {showPopup && (
        <div
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0,0,0,0.5)',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            zIndex: 1000,
          }}
        >
          <div
            style={{
              background: '#1e1e1e',
              color: '#fff',
              padding: 30,
              borderRadius: 8,
              display: 'flex',
              flexDirection: 'column',
              gap: '10px',
              position: 'relative',
              minWidth: 300,
            }}
          >
            {/* Close Button */}
            <button
              onClick={() => setShowPopup(false)}
              style={{
                position: 'absolute',
                top: 10,
                right: 10,
                background: 'transparent',
                border: 'none',
                color: '#fff',
                fontSize: 18,
                cursor: 'pointer',
              }}
            >
              ✕
            </button>

            <label>
              Algorithm:
              <select
                value={algorithm}
                onChange={(e) => setAlgorithm(e.target.value)}
                style={{
                  width: '100%',
                  padding: '8px',
                  backgroundColor: '#2c2c2c',
                  color: '#fff',
                  border: '1px solid #555',
                  borderRadius: 4,
                  marginTop: '5px',
                }}
              >
                <option value="">Select...</option>
                <option value="DFS">DFS</option>
                <option value="BFS">BFS</option>
              </select>
            </label>

            <label>
              Mode:
              <select
                value={mode}
                onChange={(e) => setMode(e.target.value)}
                style={{
                  width: '100%',
                  padding: '8px',
                  backgroundColor: '#2c2c2c',
                  color: '#fff',
                  border: '1px solid #555',
                  borderRadius: 4,
                  marginTop: '5px',
                }}
              >
                <option value="">Select...</option>
                <option value="One Solution">One Solution</option>
                <option value="All Solutions">All Solutions</option>
              </select>
            </label>

            <button
              onClick={handleFinalSearch}
              style={{
                background: '#4caf50',
                color: '#fff',
                padding: '10px',
                border: 'none',
                borderRadius: 4,
                cursor: 'pointer',
                marginTop: '10px',
              }}
            >
              Final Search
            </button>
          </div>
        </div>
      )}

      {diagramData && (
        <div style={{ flex: 1, marginTop: '20px' }}>
          <FlowDiagram data={diagramData} />
        </div>
      )}
    </main>
  );
}
