import React, { useState } from 'react';

export default function KeyValueForm() {
  const [key, setKey] = useState('');
  const [value, setValue] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    const response = await fetch('http://localhost:8080/keys', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ key, value }),
    });
    alert(response.ok ? 'Saved successfully!' : 'Error!');
  };

  return (
    <div className="form">
      <input
        type="text"
        placeholder="Key"
        value={key}
        onChange={(e) => setKey(e.target.value)}
      />
      <input
        type="text"
        placeholder="Value"
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
      <button onClick={handleSubmit}>Save</button>
    </div>
  );
}