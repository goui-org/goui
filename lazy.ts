import { useEffect, useState } from './hooks.js';
import { Component, SlotElem, component } from './elem.js';

export interface LazyConfig<T> {
    import: () => Promise<{ default: Component<T> }>
    fallback: SlotElem
}

type Lazy = {
    (config: LazyConfig<void>): Component<void>
    <T extends { [key: string]: any }>(config: LazyConfig<T>): Component<T>
};

export let lazy: Lazy = <T>(config: LazyConfig<T>): Component<T> => {
    let comp: any = undefined;
    return (props: any) => {
        let [slot, setSlot] = useState({ val: comp ?? (() => config.fallback) });
        useEffect(() => {
            if (!comp) {
                config.import().then(e => {
                    setSlot({ val: e.default });
                    comp = e.default;
                });
            }
        }, []);
        return component(slot.val, props);
    };
};
