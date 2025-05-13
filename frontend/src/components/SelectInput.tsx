'use client';

import Select from 'react-select';
import type { SingleValue } from 'react-select';

interface Option {
    value: string;
    label: string;
}

import options from './labels.json';

export default function SelectInput({ onSelect }: { onSelect: (val: string) => void }) {
    return (
        <div style={{ width: 300 }}>
            <Select<Option, false>
                options={options}
                onChange={(e: SingleValue<Option>) => e && onSelect(e.value)}
                placeholder="Select an element..."
                isSearchable
            />
        </div>
    );
}
