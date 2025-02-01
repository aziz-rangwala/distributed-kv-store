import React, { useState } from 'react';

export default function SearchKey() {
  const [key, setKey] = useState('');
  const [value, setValue] = useState('');

  const handleSearch = async () => {
    const response = await fetch(`http://localhost:8080/keys/${key}`);
    const data = await response.json();
    setValue(data.value || data.error);
  };

  return (
    <div className="search">
      <input
        type="text"
        placeholder="Search Key"
        value={key}
        onChange={(e) => setKey(e.target.value)}
      />
      <button onClick={handleSearch}>Search</button>
      <p>{value}</p>
    </div>
  );
}