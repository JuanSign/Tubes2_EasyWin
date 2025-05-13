'use client';

import dynamic from 'next/dynamic';
import { useState } from 'react';

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

  return (
    <main style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
      <form
        onSubmit={handleFirstSubmit}
        style={{ display: 'flex', gap: '10px', justifyContent: 'center', marginTop: 20 }}
      >
        <SelectInput onSelect={setSelected} />
        <button type="submit">Search</button>
      </form>

      {showPopup && (
        <div style={{
          position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.5)',
          display: 'flex', justifyContent: 'center', alignItems: 'center',
          zIndex: 1000,
        }}>
          <div style={{ background: 'white', padding: 30, borderRadius: 8, display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <label>
              Algorithm:
              <select value={algorithm} onChange={(e) => setAlgorithm(e.target.value)}>
                <option value="">Select...</option>
                <option value="DFS">DFS</option>
                <option value="BFS">BFS</option>
              </select>
            </label>

            <label>
              Mode:
              <select value={mode} onChange={(e) => setMode(e.target.value)}>
                <option value="">Select...</option>
                <option value="One Solution">One Solution</option>
                <option value="All Solutions">All Solutions</option>
              </select>
            </label>

            <button onClick={handleFinalSearch}>Final Search</button>
          </div>
        </div>
      )}

      {diagramData && (
        <div style={{ flex: 1 }}>
          <FlowDiagram data={diagramData} />
        </div>
      )}
    </main>
  );
}
