'use client';

import Select from 'react-select';
import type { SingleValue, StylesConfig } from 'react-select';

interface Option {
    value: string;
    label: string;
}

import options from './labels.json';

const customStyles: StylesConfig<Option, false> = {
    option: (provided, state) => ({
        ...provided,
        color: 'white',
        backgroundColor: state.isFocused ? '#1976d2' : '#0d47a1',
        cursor: 'pointer',
    }),
    control: (provided) => ({
        ...provided,
        backgroundColor: '#1e1e1e',
        borderColor: '#333',
        color: 'white',
    }),
    singleValue: (provided) => ({
        ...provided,
        color: 'white',
    }),
    menu: (provided) => ({
        ...provided,
        backgroundColor: '#0d47a1',
    }),
};

export default function SelectInput({ onSelect }: { onSelect: (val: string) => void }) {
    return (
        <div style={{ width: 300 }}>
            <Select<Option, false>
                options={options}
                onChange={(e: SingleValue<Option>) => e && onSelect(e.value)}
                placeholder="Select an element..."
                isSearchable
                styles={customStyles}
            />
        </div>
    );
}
